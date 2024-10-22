package gcsbackend

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Iam struct {
	pulumi.ResourceState
	name string
}

func setupIam(
	ctx *pulumi.Context,
	name string,
	args *GcsBackendArgs,
	opts ...pulumi.ResourceOption,
) (*Iam, error) {
	i := &Iam{name: name}
	if err := ctx.RegisterComponentResource(iam.String(), name, i, opts...); err != nil {
		return nil, err
	}

	for _, member := range args.IamGcsAdmins {
		if _, err := projects.NewIAMMember(ctx, fmt.Sprintf("%v/%v/%v", name, "storage.admin", member),
			&projects.IAMMemberArgs{
				Project: args.ProjectId,
				Role:    pulumi.String("roles/storage.admin"),
				Member:  pulumi.String(member),
			},
			pulumi.Parent(i),
		); err != nil {
			return nil, err
		}
	}

	var projectId *string
	args.ProjectId.ToStringOutput().ApplyT(func(id string) error {
		projectId = &id
		return nil
	})

	gcsSa, err := storage.GetProjectServiceAccount(ctx,
		&storage.GetProjectServiceAccountArgs{
			Project: projectId,
		},
		pulumi.Parent(i),
	)
	if err != nil {
		return nil, err
	}

	gcsIam, err := projects.NewIAMMember(ctx, fmt.Sprintf("%v/%v/%v", name, "cloudkms.cryptoKeyEncrypterDecrypter", gcsSa.Member),
		&projects.IAMMemberArgs{
			Project: args.ProjectId,
			Role:    pulumi.String("roles/cloudkms.cryptoKeyEncrypterDecrypter"),
			Member:  pulumi.String(gcsSa.Member),
		},
		pulumi.Parent(i),
	)
	if err != nil {
		return nil, err
	}

	err = ctx.RegisterResourceOutputs(i, pulumi.Map{
		"gcsIam": gcsIam,
	})
	if err != nil {
		return nil, err
	}

	return i, nil
}
