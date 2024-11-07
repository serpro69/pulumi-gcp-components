package project

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/serpro69/pulumi-google-components/firebase/project/util"
	"github.com/serpro69/pulumi-google-components/firebase/project/vars"
)

// FirebaseProjectIAM is a struct that represents IAM of a given firebase project
type FirebaseProjectIam struct {
	pulumi.ResourceState
}

func configureIAM(
	ctx *pulumi.Context,
	name string,
	projectId string,
	projectNumber string,
	args *vars.ProjectIamArgs,
	opts ...pulumi.ResourceOption,
) (*FirebaseProjectIam, error) {
	fpIam := &FirebaseProjectIam{}
	if err := ctx.RegisterComponentResource(util.Iam.String(), name, fpIam, opts...); err != nil {
		return nil, err
	}

	dsa, err := compute.GetDefaultServiceAccount(ctx,
		&compute.GetDefaultServiceAccountArgs{
			Project: pulumi.StringRef(projectId),
		},
		pulumi.Parent(fpIam),
	)
	if err != nil {
		return nil, err
	}

	args.ComputeServiceAccountRoles.ToStringArrayOutput().ApplyT(func(roles []string) error {
		for _, role := range roles {
			_, err := projects.NewIAMMember(ctx, fmt.Sprintf("%v/%v/%v", name, role, dsa.Member),
				&projects.IAMMemberArgs{
					Project: pulumi.String(dsa.Project),
					Role:    pulumi.String(fmt.Sprintf("roles/%v", role)),
					Member:  pulumi.String(dsa.Member),
				},
				pulumi.Parent(fpIam),
			)
			if err != nil {
				return err
			}
		}
		return nil
	})

	args.PubSubServiceAccountRoles.ToStringArrayOutput().ApplyT(func(roles []string) error {
		m := fmt.Sprintf("serviceAccount:service-%s@gcp-sa-pubsub.iam.gserviceaccount.com", projectNumber)
		for _, role := range roles {
			_, err := projects.NewIAMMember(ctx, fmt.Sprintf("%v/%v/%v", name, role, m),
				&projects.IAMMemberArgs{
					Project: pulumi.String(dsa.Project),
					Role:    pulumi.String(fmt.Sprintf("roles/%v", role)),
					Member:  pulumi.String(m),
				},
				pulumi.Parent(fpIam),
			)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err := ctx.RegisterResourceOutputs(fpIam, pulumi.Map{
		"defaultComputeSA": pulumi.String(dsa.Member),
	}); err != nil {
		return nil, err
	}

	return fpIam, nil
}
