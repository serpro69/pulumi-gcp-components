package vars

import (
	"errors"
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

	IdpArgs *identityplatform.ConfigArgs
}

// GetProjectServicesArgs returns the ProjectIdpArgs() for the FirebaseProject.
func (pa *FirebaseAuthArgs) GetIdpArgs() (*identityplatform.ConfigArgs, error) {
	args := pa.IdpArgs
	if args.Project == nil && pa.ProjectId == nil {
		return nil, errors.New("ProjectId must be set via FirebaseAuthArgs or FirebaseAuthArgs.IdpArgs")
	} else if args.Project == nil {
		args.Project = pa.ProjectId
	} else {
		pa.ProjectId = args.Project.ToStringPtrOutput().ApplyT(func(project *string) string {
			return *project
		}).(pulumi.StringInput)
	}

	authorizedDomains := make(pulumi.StringArray, 0)
	authorizedDomains = append(authorizedDomains, pulumi.String("localhost"))
	authorizedDomains = append(authorizedDomains, pulumi.Sprintf("%s.web.app", args.Project))
	authorizedDomains = append(authorizedDomains, pulumi.Sprintf("%s.firebaseapp.com", args.Project))

	if args.AuthorizedDomains == nil {
		args.AuthorizedDomains = pulumi.StringArray{}
	}
	args.AuthorizedDomains.ToStringArrayOutput().ApplyT(func(domains []string) error {
		authorizedDomains = append(authorizedDomains, pulumi.ToStringArray(domains)...)
		return nil
	})

	pa.AuthorizedApps.ToWebAppMapOutput().ApplyT(func(apps map[string]*firebase.WebApp) error {
		for name := range apps {
			log.Printf("apps: %v", name)
			authorizedDomains = append(authorizedDomains, pulumi.Sprintf("%s-%s.web.app", name, args.Project))
			authorizedDomains = append(authorizedDomains, pulumi.Sprintf("%s-%s.firebaseapp.com", name, args.Project))
		}
		return nil
	})

	args.AuthorizedDomains = authorizedDomains.ToStringArrayOutput().ApplyT(func(domains []string) []string {
		return utils.Unique(domains)
	}).(pulumi.StringArrayInput)

	return args, nil
}
