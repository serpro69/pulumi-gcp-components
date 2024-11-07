package vars

import (
	p "github.com/serpro69/pulumi-google-components/project/vars"
)

type ProjectArgs struct {
	p.ProjectArgs
	*ProjectIamArgs
}

type ProjectIamArgs struct {
	ComputeServiceAccountRoles []string `pulumi:"computeServiceAccountRoles"`
	PubSubServiceAccountRoles  []string `pulumi:"pubSubServiceAccountRoles"`
	FirebaseAdminMembers       []string `pulumi:"firebaseAdminMembers"`
	FirebaseViewerMembers      []string `pulumi:"firebaseViewerMembers"`
}

func DefaultProjectIamArgs() *ProjectIamArgs {
	return &ProjectIamArgs{
		// firebase deployment-/runtime-related roles for default compute SA
		ComputeServiceAccountRoles: []string{
			"artifactregistry.createOnPushWriter",
			"eventarc.eventReceiver",
			"firebase.admin",
			"logging.logWriter",
			"run.invoker",
			"serviceusage.serviceUsageConsumer",
			"storage.objectViewer",
		},
	}
}
