package project

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/serpro69/pulumi-google-components/firebase/project/util"
	"github.com/serpro69/pulumi-google-components/firebase/project/vars"
	"github.com/serpro69/pulumi-google-components/project"
	"github.com/serpro69/pulumi-google-components/utils"
)

// FirebaseProject is a struct that represents a project with enabled Firebase support in GCP
type FirebaseProject struct {
	pulumi.ResourceState

	*project.Project
}

// NewProject creates a new Project in GCP
func NewFirebaseProject(
	ctx *pulumi.Context,
	name string,
	args *vars.ProjectArgs,
	opts ...pulumi.ResourceOption,
) (*FirebaseProject, error) {
	p := &FirebaseProject{}
	err := ctx.RegisterComponentResource(util.Project.String(), name, p, opts...)
	if err != nil {
		return nil, err
	}

	// Required for the project to display in any list of Firebase projects.
	args.Labels = args.Labels.ToStringMapOutput().ApplyT(func(labels map[string]string) map[string]string {
		if l, found := labels["firebase"]; !found || l != "enabled" {
			labels["firebase"] = "enabled"
		}
		return labels
	}).(pulumi.StringMapInput)

	args.ActivateApis = args.ActivateApis.ToStringArrayOutput().ApplyT(func(apis []string) []string {
		apis = append(apis, services...)
		return utils.Unique(apis)
	}).(pulumi.StringArrayInput)

	if p.Project, err = project.NewProject(ctx, name, &args.ProjectArgs, pulumi.Parent(p)); err != nil {
		return nil, err
	}

	// add proj to wait until project is fully created (and services are enabled)
	pulumi.All(p.Project.Main.ProjectId, p.Project.Main.Number, p.Project).ApplyT(func(all []interface{}) error {
		projectId := all[0].(string)
		projectNumber := all[1].(string)
		configureIAM(ctx, name, projectId, projectNumber, args.GetProjectIamArgs(), pulumi.Parent(p), pulumi.DeletedWith(p.Project))
		return nil
	})

	if err := ctx.RegisterResourceOutputs(p, pulumi.Map{
		"projectId": p.Project.Main.ProjectId.ToStringOutput(),
	}); err != nil {
		return nil, err
	}

	return p, nil
}

var services = []string{
	// base
	"cloudbilling.googleapis.com",
	"cloudresourcemanager.googleapis.com",
	// By enabling the Service Usage API, the project will be able to accept quota checks!
	// So, for all subsequent resource provisioning and service enabling, you should use the provider with user_project_override (no alias needed).
	"serviceusage.googleapis.com",
	// firebase services
	"firebase.googleapis.com",
	"fcm.googleapis.com",
	"fcmregistrations.googleapis.com",
	"firebaseappdistribution.googleapis.com",
	"firebaseextensions.googleapis.com",
	"firebasedynamiclinks.googleapis.com",
	"firebasehosting.googleapis.com",
	"firebaseinstallations.googleapis.com",
	"firebaseremoteconfig.googleapis.com",
	"firebaseremoteconfigrealtime.googleapis.com",
	"firebaserules.googleapis.com",
	// firebase functions
	// i  functions: ensuring required API cloudfunctions.googleapis.com is enabled...
	// i  functions: ensuring required API cloudbuild.googleapis.com is enabled...
	// i  artifactregistry: ensuring required API artifactregistry.googleapis.com is enabled...
	"cloudfunctions.googleapis.com",
	"cloudbuild.googleapis.com",
	"artifactregistry.googleapis.com",
	// i  functions: packaged .../foo/.firebase/bar/functions (32.51 MB) for uploading
	// i  functions: ensuring required API run.googleapis.com is enabled...
	// i  functions: ensuring required API eventarc.googleapis.com is enabled...
	// i  functions: ensuring required API pubsub.googleapis.com is enabled...
	// i  functions: ensuring required API storage.googleapis.com is enabled...
	"run.googleapis.com",
	"eventarc.googleapis.com",
	"pubsub.googleapis.com",
	"storage.googleapis.com",
}
