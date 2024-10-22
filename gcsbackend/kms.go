package gcsbackend

import (
	gcpKms "github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/kms"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Kms struct {
	pulumi.ResourceState
	name      string
	KeyRing   *gcpKms.KeyRing
	CryptoKey *gcpKms.CryptoKey
}

func setupKms(
	ctx *pulumi.Context,
	name string,
	args *GcsBackendArgs,
	l *locals,
	opts ...pulumi.ResourceOption,
) (*Kms, error) {
	k := &Kms{name: name}
	if err := ctx.RegisterComponentResource(kms.String(), name, k, opts...); err != nil {
		return nil, err
	}

	kr, err := gcpKms.NewKeyRing(ctx, name,
		&gcpKms.KeyRingArgs{
			Project:  args.ProjectId,
			Name:     l.kmsKeyRingName,
			Location: l.kmsLocation,
		},
		pulumi.Parent(k),
	)
	if err != nil {
		return nil, err
	} else {
		k.KeyRing = kr
	}

	dsd := pulumi.String("2592000s")
	protect := pulumi.Protect(false)

	if args.IsProd.First() {
		dsd = pulumi.String("7776000s")
		protect = pulumi.Protect(true)
	}

	kc, err := gcpKms.NewCryptoKey(ctx, name,
		&gcpKms.CryptoKeyArgs{
			Name:                     l.kmsKeyName,
			KeyRing:                  k.KeyRing.ID(),
			RotationPeriod:           pulumi.String("2592000s"), // 30 days
			DestroyScheduledDuration: dsd,
			Purpose:                  pulumi.String("ENCRYPT_DECRYPT"),
		},
		pulumi.Parent(k),
		protect,
	)
	if err != nil {
		return nil, err
	} else {
		k.CryptoKey = kc
	}

	err = ctx.RegisterResourceOutputs(k,
		pulumi.Map{
			"keyRing":   kr,
			"cryptoKey": kc,
		})
	if err != nil {
		return nil, err
	}

	return k, nil
}

// func newManagedFolder(ctx *pulumi.Context, name string, b *storage.Bucket) (*storage.ManagedFolder, error) {
// 	// managed folder name must end with '/'
// 	if name[len(name)-1] != '/' {
// 		name = name + "/"
// 	}
// 	f, err := storage.NewManagedFolder(ctx, name,
// 		&storage.ManagedFolderArgs{
// 			Bucket: b.Name,
// 			Name:   pulumi.String(name),
// 		},
// 		pulumi.Parent(b),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return f, nil
// }
//
// func newManagedFolderIamPolicy(ctx *pulumi.Context, name string, f *storage.ManagedFolder, admins []string, users []string) error {
// 	p, err := organizations.LookupIAMPolicy(ctx,
// 		&organizations.LookupIAMPolicyArgs{
// 			Bindings: []organizations.GetIAMPolicyBinding{
// 				{
// 					Role:    "roles/storage.admin",
// 					Members: admins,
// 				},
// 				{
// 					Role:    "roles/storage.objectUser",
// 					Members: users,
// 				},
// 			},
// 		},
// 		pulumi.Parent(f),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = storage.NewManagedFolderIamPolicy(ctx, name,
// 		&storage.ManagedFolderIamPolicyArgs{
// 			Bucket:        f.Bucket,
// 			ManagedFolder: f.Name,
// 			PolicyData:    pulumi.String(p.PolicyData),
// 		},
// 		pulumi.Parent(f),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
