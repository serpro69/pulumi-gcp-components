package gcsbackend

import (
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Bucket struct {
	pulumi.ResourceState
	name    string
	Main    *storage.Bucket
	Folders []*storage.ManagedFolder
}

func newGcsBucket(
	ctx *pulumi.Context,
	name string,
	args *GcsBackendArgs,
	l *locals,
	k *Kms,
	opts ...pulumi.ResourceOption,
) (*Bucket, error) {
	b := &Bucket{name: name}
	if err := ctx.RegisterComponentResource(bucket.String(), name, b, opts...); err != nil {
		return nil, err
	}

	rd := pulumi.Int(604800)
	nvState := pulumi.Int(11)
	nvLock := pulumi.Int(0)

	if args.IsProd.First() {
		rd = pulumi.Int(7776000)
		nvState = pulumi.Int(31)
		nvLock = pulumi.Int(4)
	}

	bucket, err := storage.NewBucket(ctx, name,
		&storage.BucketArgs{
			Project:      args.ProjectId,
			Name:         l.statePrefix,
			ForceDestroy: args.IsProd.Second(),
			// https://cloud.google.com/storage/docs/locations#predefined
			Location:     pulumi.String("EUR4"),
			StorageClass: pulumi.String("STANDARD"),

			// https://cloud.google.com/storage/docs/uniform-bucket-level-access#should-you-use
			// https://cloud.google.com/storage/docs/using-uniform-bucket-level-access#command-line_1
			UniformBucketLevelAccess: pulumi.Bool(true),

			// Prevent public access irrespective of org policy
			PublicAccessPrevention: pulumi.String("enforced"),

			Versioning: &storage.BucketVersioningArgs{
				Enabled: pulumi.Bool(true),
			},

			SoftDeletePolicy: &storage.BucketSoftDeletePolicyArgs{
				RetentionDurationSeconds: rd,
			},

			Encryption: &storage.BucketEncryptionArgs{
				DefaultKmsKeyName: k.CryptoKey.ID(),
			},

			LifecycleRules: storage.BucketLifecycleRuleArray{
				&storage.BucketLifecycleRuleArgs{
					Action: &storage.BucketLifecycleRuleActionArgs{
						Type: pulumi.String("Delete"),
					},
					Condition: &storage.BucketLifecycleRuleConditionArgs{
						MatchesSuffixes: pulumi.StringArray{
							pulumi.String(".tfstate"),
						},
						WithState: pulumi.String("ARCHIVED"),
						MatchesStorageClasses: pulumi.StringArray{
							pulumi.String("STANDARD"),
						},
						NumNewerVersions: nvState,
					},
				},
				&storage.BucketLifecycleRuleArgs{
					Action: &storage.BucketLifecycleRuleActionArgs{
						Type: pulumi.String("Delete"),
					},
					Condition: &storage.BucketLifecycleRuleConditionArgs{
						MatchesSuffixes: pulumi.StringArray{
							pulumi.String(".tflock"),
						},
						WithState: pulumi.String("ARCHIVED"),
						MatchesStorageClasses: pulumi.StringArray{
							pulumi.String("STANDARD"),
						},
						NumNewerVersions: nvLock,
					},
				},
			},
		},
		pulumi.Parent(b),
	)
	if err != nil {
		return nil, err
	}

	b.Main = bucket
	b.Folders = make([]*storage.ManagedFolder, 0)

	if p, err := organizations.LookupIAMPolicy(ctx,
		&organizations.LookupIAMPolicyArgs{
			Bindings: []organizations.GetIAMPolicyBinding{
				{
					Role:    "roles/storage.admin",
					Members: args.IamGcsAdmins,
				},
				{
					Role:    "roles/storage.objectViewer",
					Members: args.IamGcsObjectViewers,
				},
			},
		},
		pulumi.Parent(b),
	); err != nil {
		return nil, err
	} else {
		_, err = storage.NewBucketIAMPolicy(ctx, name,
			&storage.BucketIAMPolicyArgs{
				Bucket:     b.Main.Name,
				PolicyData: pulumi.String(p.PolicyData),
			},
			pulumi.Parent(b),
		)
		if err != nil {
			return nil, err
		}
	}

	if f, e := newManagedFolder(ctx, "state/prod/.pulumi", b.Main); e != nil {
		return nil, e
	} else {
		if e := newManagedFolderIamPolicy(
			ctx,
			"state/prod/.pulumi/",
			f,
			args.IamGcsPlStateFolderProdAdmins,
			args.IamGcsPlStateFolderProdUsers,
		); e != nil {
			return nil, e
		}
		b.Folders = append(b.Folders, f)
	}

	if f, e := newManagedFolder(ctx, "state/test/.pulumi/", b.Main); e != nil {
		return nil, e
	} else {
		if e := newManagedFolderIamPolicy(
			ctx,
			"state/test/.pulumi/",
			f,
			args.IamGcsPlStateFolderTestAdmins,
			args.IamGcsPlStateFolderTestUsers,
		); e != nil {
			return nil, e
		}
		b.Folders = append(b.Folders, f)
	}

	err = ctx.RegisterResourceOutputs(b,
		pulumi.Map{
			"bucket":  bucket,
			"folders": pulumi.All(b.Folders),
		})
	if err != nil {
		return nil, err
	}

	return b, nil
}

func newManagedFolder(ctx *pulumi.Context, name string, b *storage.Bucket) (*storage.ManagedFolder, error) {
	// managed folder name must end with '/'
	if name[len(name)-1] != '/' {
		name = name + "/"
	}
	f, err := storage.NewManagedFolder(ctx, name,
		&storage.ManagedFolderArgs{
			Bucket: b.Name,
			Name:   pulumi.String(name),
		},
		pulumi.Parent(b),
	)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func newManagedFolderIamPolicy(ctx *pulumi.Context, name string, f *storage.ManagedFolder, admins []string, users []string) error {
	p, err := organizations.LookupIAMPolicy(ctx,
		&organizations.LookupIAMPolicyArgs{
			Bindings: []organizations.GetIAMPolicyBinding{
				{
					Role:    "roles/storage.admin",
					Members: admins,
				},
				{
					Role:    "roles/storage.objectUser",
					Members: users,
				},
			},
		},
		pulumi.Parent(f),
	)
	if err != nil {
		return err
	}
	_, err = storage.NewManagedFolderIamPolicy(ctx, name,
		&storage.ManagedFolderIamPolicyArgs{
			Bucket:        f.Bucket,
			ManagedFolder: f.Name,
			PolicyData:    pulumi.String(p.PolicyData),
		},
		pulumi.Parent(f),
	)
	if err != nil {
		return err
	}
	return nil
}
