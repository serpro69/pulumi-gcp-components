package project

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-time/sdk/go/time"
	"github.com/serpro69/pulumi-google-components/project/services"
	"github.com/serpro69/pulumi-google-components/project/util"
	"github.com/serpro69/pulumi-google-components/project/vars"
)

// Project is a struct that represents a project in GCP
type Project struct {
	pulumi.ResourceState

	// TODO: can't use composition via `*organizations.Project` because it's in turn composed of `pulumi.CustomResourceState`, which contains `getProviders()` method that needs to be implemented?

	Main *organizations.Project `pulumi:"project"`
	*services.ProjectServices
}

// NewProject creates a new Project in GCP
func NewProject(
	ctx *pulumi.Context,
	name string,
	args *vars.ProjectArgs,
	opts ...pulumi.ResourceOption,
) (*Project, error) {
	p := &Project{}
	err := ctx.RegisterComponentResource(util.Project.String(), name, p, opts...)
	if err != nil {
		return nil, err
	}
	p.Main, err = organizations.NewProject(ctx, name,
		&organizations.ProjectArgs{
			BillingAccount:    args.BillingAccount,
			FolderId:          args.FolderId,
			ProjectId:         args.ProjectId,
			Name:              args.ProjectName,
			AutoCreateNetwork: args.AutoCreateNetwork,
			Labels:            args.Labels,
			DeletionPolicy:    args.DeletionPolicy,
		},
		pulumi.Parent(p),
	)
	if err != nil {
		return nil, err
	}

	// https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/google_project_service#mitigation---adding-sleeps
	wfp, err := time.NewSleep(ctx, fmt.Sprintf("wait"),
		&time.SleepArgs{
			CreateDuration: pulumi.String("30s"),
		},
		pulumi.Parent(p),
		pulumi.DeletedWith(p.Main),
		pulumi.DependsOn([]pulumi.Resource{p.Main}),
	)
	if err != nil {
		return nil, err
	}

	// Activate Services
	p.ProjectServices, err = services.ActivateApis(ctx, name, args.GetProjectServicesArgs(),
		pulumi.Parent(p),
		pulumi.DependsOn([]pulumi.Resource{wfp}),
		pulumi.DeletedWith(p.Main),
	)
	if err != nil {
		return nil, err
	}

	// Create IAM members
	if _, err := newIamMember(ctx, p, "owner", args.Owners,
		pulumi.DependsOn([]pulumi.Resource{wfp}),
	); err != nil {
		return nil, err
	}
	if _, err := newIamMember(ctx, p, "editor", args.Editors,
		pulumi.DependsOn([]pulumi.Resource{wfp}),
	); err != nil {
		return nil, err
	}
	if _, err := newIamMember(ctx, p, "viewer", args.Viewers,
		pulumi.DependsOn([]pulumi.Resource{wfp}),
	); err != nil {
		return nil, err
	}

	err = ctx.RegisterResourceOutputs(p,
		pulumi.Map{
			"main": p.Main,
			"wait": wfp,
			"apis": p.Services,
		})
	if err != nil {
		return nil, err
	}
	return p, nil
}

// newIamMember creates a list of IAM members in a GCP Project with a given role
func newIamMember(ctx *pulumi.Context, parent *Project, role string, members pulumi.StringArray, opts ...pulumi.ResourceOption) ([]*projects.IAMMember, error) {
	mm := []*projects.IAMMember{}
	for _, m := range members {
		if res, err := projects.NewIAMMember(ctx, fmt.Sprintf("%v$%v", m, role),
			&projects.IAMMemberArgs{
				Project: parent.Main.ProjectId,
				Role:    pulumi.String(fmt.Sprintf("roles/%s", role)),
				Member:  m,
			},
			append(opts, pulumi.Parent(parent), pulumi.DeletedWith(parent.Main))...,
		); err != nil {
			return nil, err
		} else {
			mm = append(mm, res)
		}
	}
	return mm, nil
}
