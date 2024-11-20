package vars

import (
	"slices"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	p "github.com/serpro69/pulumi-google-components/project/vars"
	"github.com/serpro69/pulumi-google-components/utils"
)

type ProjectArgs struct {
	p.ProjectArgs
	*projectIamArgs
}

func DefaultProjectArgs() *ProjectArgs {
	dpa := p.DefaultProjectArgs()

	return &ProjectArgs{
		ProjectArgs:    *dpa,
		projectIamArgs: defaultProjectIamArgs(),
	}
}

func (pa *ProjectArgs) GetProjectIamArgs() *ProjectIamArgs {
	args := pa.projectIamArgs

	pa.GetProjectServicesArgs().ActivateApis.ToStringArrayOutput().ApplyT(func(apis []string) error {
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
		args.FirebaseAdminMembers = pa.ProjectArgs.Owners
	}
	if len(args.FirebaseViewerMembers) == 0 {
		args.FirebaseViewerMembers = pa.ProjectArgs.Viewers
	}

	return &ProjectIamArgs{args}
}

type projectIamArgs struct {
	ComputeServiceAccountRoles pulumi.StringArray
	PubSubServiceAccountRoles  pulumi.StringArray
	FirebaseAdminMembers       pulumi.StringArray
	FirebaseViewerMembers      pulumi.StringArray
}

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
