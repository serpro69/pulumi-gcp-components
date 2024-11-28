package auth

import (
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/identityplatform"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/serpro69/pulumi-google-components/firebase/auth/util"
	"github.com/serpro69/pulumi-google-components/firebase/auth/vars"
	"github.com/serpro69/pulumi-google-components/project/services"
	sv "github.com/serpro69/pulumi-google-components/project/vars"
)

// FirebaseProject is a struct that represents a project with enabled Firebase support in GCP
type FirebaseAuth struct {
	pulumi.ResourceState

	IdpConfig *identityplatform.Config
}

// NewFirebaseAuthConfig configures firease auth for a project
func NewFirebaseAuthConfig(
	ctx *pulumi.Context,
	name string,
	args *vars.FirebaseAuthArgs,
	opts ...pulumi.ResourceOption,
) (*FirebaseAuth, error) {
	a := &FirebaseAuth{}
	err := ctx.RegisterComponentResource(util.Idp.String(), name, a, opts...)
	if err != nil {
		return nil, err
	}

	idpArgs, err := args.GetIdpArgs()
	if err != nil {
		return nil, err
	}

	ss, err := services.ActivateApis(ctx, name,
		&sv.ProjectServicesArgs{
			ProjectId: idpArgs.Project.ToStringPtrOutput().ApplyT(
				func(project *string) string { return *project },
			).(pulumi.StringInput),
			ActivateApis: pulumi.ToStringArray(apis),
		},
		pulumi.Parent(a),
	)
	if err != nil {
		return nil, err
	}

	if a.IdpConfig, err = setupIdp(ctx, name, idpArgs,
		pulumi.Parent(a),
		pulumi.DependsOn([]pulumi.Resource{ss}),
	); err != nil {
		return nil, err
	}

	if err := ctx.RegisterResourceOutputs(a, pulumi.Map{
		"idp": a.IdpConfig,
	}); err != nil {
		return nil, err
	}

	return a, nil
}

var apis = []string{
	"cloudbilling.googleapis.com",
	"cloudresourcemanager.googleapis.com",
	"serviceusage.googleapis.com",
	"identitytoolkit.googleapis.com",
}

func setupIdp(ctx *pulumi.Context, name string, args *identityplatform.ConfigArgs, opts ...pulumi.ResourceOption) (*identityplatform.Config, error) {
	idp, err := identityplatform.NewConfig(ctx, name, args, opts...)
	if err != nil {
		return nil, err
	}
	return idp, nil
}
