package vars

import (
	"slices"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pv "github.com/serpro69/pulumi-google-components/project/vars"
	"github.com/serpro69/pulumi-google-components/utils"
)

type ProjectArgs struct {
	Project *pv.ProjectArgs // avoid shadowing of ProjectId
	*projectIamArgs
	*projectWebAppsArgs
}

// DefaultProjectArgs returns a default set of arguments for a FirebaseProject.
func DefaultProjectArgs(pa *pv.ProjectArgs) *ProjectArgs {
	return &ProjectArgs{
		Project:            pa,
		projectIamArgs:     defaultProjectIamArgs(),
		projectWebAppsArgs: defaultProjectWebAppsArgs(pa.ProjectId),
	}
}

// GetProjectServicesArgs returns the ProjectServicesArgs for the FirebaseProject.
func (pa *ProjectArgs) GetProjectIamArgs() *ProjectIamArgs {
	args := pa.projectIamArgs

	pa.Project.GetProjectServicesArgs().ActivateApis.ToStringArrayOutput().ApplyT(func(apis []string) error {
		if slices.Contains(apis, "cloudfunctions.googleapis.com") {
			// firebase functions deployment-/runtime-related roles for default compute SA
			args.ComputeServiceAccountRoles = utils.Unique(
				append(args.ComputeServiceAccountRoles,
					pulumi.ToStringArray([]string{
						"artifactregistry.createOnPushWriter",
						"eventarc.eventReceiver",
						"firebase.admin",
						"logging.logWriter",
						"run.invoker",
						"serviceusage.serviceUsageConsumer",
						"storage.objectViewer",
					})...,
				),
			)
		}
		return nil
	})

	if len(args.FirebaseAdminMembers) == 0 {
		args.FirebaseAdminMembers = pa.Project.Owners
	}
	if len(args.FirebaseViewerMembers) == 0 {
		args.FirebaseViewerMembers = pa.Project.Viewers
	}

	return &ProjectIamArgs{args}
}

type projectIamArgs struct {
	ComputeServiceAccountRoles pulumi.StringArray
	PubSubServiceAccountRoles  pulumi.StringArray
	FirebaseAdminMembers       pulumi.StringArray
	FirebaseViewerMembers      pulumi.StringArray
}

// ProjectIamArgs represents the arguments for configuring IAM roles for a FirebaseProject.
type ProjectIamArgs struct {
	*projectIamArgs
}

func defaultProjectIamArgs() *projectIamArgs {
	return &projectIamArgs{
		ComputeServiceAccountRoles: pulumi.ToStringArray(make([]string, 0)),
		PubSubServiceAccountRoles:  pulumi.ToStringArray(make([]string, 0)),
		FirebaseAdminMembers:       pulumi.ToStringArray(make([]string, 0)),
		FirebaseViewerMembers:      pulumi.ToStringArray(make([]string, 0)),
	}
}

// GetProjectWebAppsArgs returns the ProjectWebAppsArgs for the FirebaseProject.
func (pa *ProjectArgs) GetProjectWebAppsArgs() *ProjectWebAppsArgs {
	return &ProjectWebAppsArgs{pa.projectWebAppsArgs}
}

type projectWebAppsArgs struct {
	// Project ID to enable APIs on.
	// Mandatory value. An error will be returned if ProjectId is not set.
	ProjectId pulumi.StringInput
	// The list of web apps to create within the project
	WebApps pulumi.StringArray
	// A map of WebApps associated with a list of custom domains
	CustomDomains pulumi.StringArrayMap
}

func defaultProjectWebAppsArgs(projectId pulumi.StringInput) *projectWebAppsArgs {
	return &projectWebAppsArgs{
		ProjectId:     projectId,
		WebApps:       pulumi.ToStringArray(make([]string, 0)),
		CustomDomains: pulumi.ToStringArrayMap(make(map[string][]string)),
	}
}

// ProjectWebAppsArgs represents the arguments for configuring WebApps for a FirebaseProject.
type ProjectWebAppsArgs struct {
	*projectWebAppsArgs
}
