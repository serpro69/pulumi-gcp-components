package vars

import (
	"log"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/firebase"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/identityplatform"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/serpro69/pulumi-google-components/utils"
)

type FirebaseAuthArgs struct {
	ProjectId pulumi.StringInput `pulumi:"projectId"`

	// List of firebase applications to add to authorized domains
	AuthorizedApps firebase.WebAppMap `pulumi:"authorizedApps"`

	*projectIdpArgs
}

// DefaultFirebaseAuthArgs returns a default set of arguments for a FirebaseAuth
func DefaultFirebaseAuthArgs() *FirebaseAuthArgs {
	return &FirebaseAuthArgs{
		AuthorizedApps: make(firebase.WebAppMap),
		projectIdpArgs: defaultProjectIdpArgs(),
	}
}

// GetProjectServicesArgs returns the ProjectIdpArgs() for the FirebaseProject.
func (pa *FirebaseAuthArgs) GetProjectIdpArgs() *ProjectIdpArgs {
	args := pa.projectIdpArgs
	authorizedDomains := make(pulumi.StringArray, 0)
	authorizedDomains = append(authorizedDomains, pulumi.String("localhost"))
	authorizedDomains = append(authorizedDomains, pulumi.Sprintf("%s.web.app", pa.ProjectId))
	authorizedDomains = append(authorizedDomains, pulumi.Sprintf("%s.firebaseapp.com", pa.ProjectId))
	authorizedDomains = append(authorizedDomains, pa.AuthorizedDomains...)

	pa.AuthorizedApps.ToWebAppMapOutput().ApplyT(func(apps map[string]*firebase.WebApp) error {
		for name := range apps {
			log.Printf("apps: %v", name)
			authorizedDomains = append(authorizedDomains, pulumi.Sprintf("%s-%s.web.app", name, pa.ProjectId))
			authorizedDomains = append(authorizedDomains, pulumi.Sprintf("%s-%s.firebaseapp.com", name, pa.ProjectId))
		}
		return nil
	})

	args.AuthorizedDomains = utils.Unique(authorizedDomains)
	return &ProjectIdpArgs{args}
}

type projectIdpArgs struct {
	// Additional authorized domains for firebase hosting
	AuthorizedDomains pulumi.StringArray `pulumi:"authorizedDomains"`

	// OIDC client secret
	OidcClientSecret pulumi.StringInput `pulumi:"oidcClientSecret"`

	// OIDC client id
	OidcClientId pulumi.StringInput `pulumi:"oidcClientId"`

	// OIDC issuer URI
	OidcIssuerUri pulumi.StringInput `pulumi:"oidcIssuerUri"`

	SignIn *identityplatform.ConfigSignIn
}

// ProjectIamArgs represents the arguments for configuring IAM roles for a FirebaseProject.
type ProjectIdpArgs struct {
	*projectIdpArgs
}

func defaultProjectIdpArgs() *projectIdpArgs {
	return &projectIdpArgs{
		AuthorizedDomains: make(pulumi.StringArray, 0),
		SignIn: &identityplatform.ConfigSignIn{
			AllowDuplicateEmails: pulumi.BoolRef(false),
			Email: &identityplatform.ConfigSignInEmail{
				Enabled:          true,
				PasswordRequired: pulumi.BoolRef(false),
			},
		},
	}
}
