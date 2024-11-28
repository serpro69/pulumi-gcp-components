package auth

import (
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/identityplatform"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/serpro69/pulumi-google-components/firebase/auth/util"
	"github.com/serpro69/pulumi-google-components/firebase/auth/vars"
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

	if a.IdpConfig, err = setupIdp(ctx, name, args.GetProjectIdpArgs(), pulumi.Parent(a)); err != nil {
		return nil, err
	}

	if err := ctx.RegisterResourceOutputs(a, pulumi.Map{
		"idp": a.IdpConfig,
	}); err != nil {
		return nil, err
	}

	return a, nil
}

var services = []string{
	"cloudbilling.googleapis.com",
	"cloudresourcemanager.googleapis.com",
	"serviceusage.googleapis.com",
	"identitytoolkit.googleapis.com",
}

func setupIdp(ctx *pulumi.Context, name string, args *vars.ProjectIdpArgs, opts ...pulumi.ResourceOption) (*identityplatform.Config, error) {
	signIn := &identityplatform.ConfigSignInArgs{}
	if args.SignIn != nil {
		if args.SignIn.Anonymous != nil {
			signIn.Anonymous = &identityplatform.ConfigSignInAnonymousArgs{
				Enabled: pulumi.Bool(args.SignIn.Anonymous.Enabled),
			}
		}
		if args.SignIn.Email != nil {
			signIn.Email = &identityplatform.ConfigSignInEmailArgs{
				Enabled:          pulumi.Bool(args.SignIn.Email.Enabled),
				PasswordRequired: pulumi.Bool(*args.SignIn.Email.PasswordRequired),
			}
		}
		if args.SignIn.PhoneNumber != nil {
			signIn.PhoneNumber = &identityplatform.ConfigSignInPhoneNumberArgs{
				Enabled:          pulumi.Bool(args.SignIn.PhoneNumber.Enabled),
				TestPhoneNumbers: pulumi.ToStringMap(args.SignIn.PhoneNumber.TestPhoneNumbers),
			}
		}
		if args.SignIn.AllowDuplicateEmails != nil {
			signIn.AllowDuplicateEmails = pulumi.Bool(*args.SignIn.AllowDuplicateEmails)
		}
	}
	c := &identityplatform.ConfigArgs{
		AuthorizedDomains: args.AuthorizedDomains,
		SignIn:            signIn,
		Client: &identityplatform.ConfigClientArgs{
			Permissions: identityplatform.ConfigClientPermissionsArgs{
				DisabledUserSignup:   pulumi.Bool(true), // TODO: make this configurable
				DisabledUserDeletion: pulumi.Bool(true), // TODO: make this configurable
			},
		},
	}
	idp, err := identityplatform.NewConfig(ctx, name, c, opts...)
	if err != nil {
		return nil, err
	}
	return idp, nil
}
