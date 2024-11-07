package vars

import (
	"slices"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ProjectArgs struct {
	// GCP billing account ID
	BillingAccount pulumi.StringInput
	// GCP folder ID
	FolderId pulumi.StringInput
	// GCP project name
	ProjectName pulumi.StringInput
	// The deletion policy for the Project.
	// Setting PREVENT will protect the project against any destroy actions caused by a terraform apply or terraform destroy.
	// Setting ABANDON allows the resource to be abandoned rather than deleted.
	// Possible values are: "PREVENT", "ABANDON", "DELETE"
	DeletionPolicy pulumi.StringInput
	// Create the 'default' network automatically
	AutoCreateNetwork pulumi.BoolInput
	// A set of key/value label pairs to assign to the project
	Labels pulumi.StringMapInput

	// Project IAM

	// Optional list of IAM-format members to set as project owners
	Owners pulumi.StringArray
	// Optional list of IAM-format members to set as project editor
	Editors pulumi.StringArray
	// Optional list of IAM-format members to set as project viewers
	Viewers pulumi.StringArray

	// Project APIs

	// Whether to disable the compute engine API if not explicitly enabled
	// It's usually not recommended to disable the compute engine API
	// as it's required for many other services
	DisableComputeEngine bool

	*projectServicesArgs
}

/*
ToProjectServicesArgs converts ProjectArgs to ProjectServicesArgs
*/
func (pa *ProjectArgs) GetProjectServicesArgs() *ProjectServicesArgs {
	psa := pa.projectServicesArgs
	if !pa.DisableComputeEngine {
		psa.ActivateApis = psa.ActivateApis.ToStringArrayOutput().ApplyT(func(apis []string) []string {
			if !slices.Contains(apis, "compute.googleapis.com") {
				return append(apis, "compute.googleapis.com")
			}
			return apis
		}).(pulumi.StringArrayInput)
	}
	return &ProjectServicesArgs{psa}
}

/*
DefaultProjectArgs returns the ProjectArgs with initialized default values:

  - AutoCreateNetwork: false
  - DisableComputeEngine: false
  - ProjectServicesArgs.DisableServicesOnDestroy: true
  - ProjectServicesArgs.DisableDependentServices: true

If DisableComputeEngine is not set to true,
the compute engine API will be added to ActivateApis list
*/
func DefaultProjectArgs() *ProjectArgs {
	return &ProjectArgs{
		AutoCreateNetwork:    pulumi.Bool(false),
		DisableComputeEngine: false,
		projectServicesArgs: &projectServicesArgs{
			DisableServicesOnDestroy: true,
			DisableDependentServices: true,
		},
	}
}

type projectServicesArgs struct {
	// Project ID to enable APIs on.
	// Mandatory value. An error will be returned if ProjectId is not set.
	ProjectId pulumi.StringInput
	// The list of apis to activate within the project
	ActivateApis pulumi.StringArrayInput
	// If `true`, disable the service when the resource is destroyed.
	// Defaults to `true`.
	// https://www.terraform.io/docs/providers/google/r/google_project_service.html#disable_on_destroy
	DisableServicesOnDestroy bool
	// If `true`, services that are enabled and which depend on this service should also be disabled when this service is destroyed.
	// If `false` or unset, an error will be generated if any enabled services depend on this service when destroying it.
	// https://www.terraform.io/docs/providers/google/r/google_project_service.html#disable_dependent_services
	DisableDependentServices bool
}

type ProjectServicesArgs struct {
	*projectServicesArgs
}
