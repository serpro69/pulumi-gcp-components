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
	// GCP project ID
	ProjectId pulumi.StringInput
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

	// The list of apis to activate within the project
	ActivateApis pulumi.StringArray
	// Whether project services will be disabled when the resources are destroyed.
	// https://www.terraform.io/docs/providers/google/r/google_project_service.html#disable_on_destroy
	DisableServicesOnDestroy bool
	// Whether services that are enabled and which depend on this service should also be disabled when this service is destroyed.
	// https://www.terraform.io/docs/providers/google/r/google_project_service.html#disable_dependent_services
	DisableDependentServices bool
	// Whether to disable the compute engine API if not explicitly enabled
	// It's usually not recommended to disable the compute engine API
	// as it's required for many other services
	DisableComputeEngine bool
}

/*
ToProjectServicesArgs converts ProjectArgs to ProjectServicesArgs
*/
func (pa *ProjectArgs) ToProjectServicesArgs() *ProjectServicesArgs {
	args := &ProjectServicesArgs{
		ProjectId:                pa.ProjectId,
		ActivateApis:             pa.ActivateApis,
		DisableOnDestroy:         pa.DisableServicesOnDestroy,
		DisableDependentServices: pa.DisableDependentServices,
	}
	if !pa.DisableComputeEngine {
		if !slices.Contains(args.ActivateApis, pulumi.StringInput(pulumi.String("compute.googleapis.com"))) {
			args.ActivateApis = append(args.ActivateApis, pulumi.String("compute.googleapis.com"))
		}
	}
	return args
}

/*
DefaultProjectArgs returns the ProjectArgs with initialized default values:

  - DisableServicesOnDestroy: true
  - DisableDependentServices: true
  - DisableComputeEngine: false

If DisableComputeEngine is not set to true,
the compute engine API will be added to ActivateApis list
*/
func DefaultProjectArgs() *ProjectArgs {
	return &ProjectArgs{
		DisableServicesOnDestroy: true,
		DisableDependentServices: true,
		DisableComputeEngine:     false,
	}
}

type ProjectServicesArgs struct {
	// Project ID to enable APIs on.
	// Mandatory value. An error will be returned if ProjectId is not set.
	ProjectId pulumi.StringInput
	// The list of apis to activate within the project
	ActivateApis pulumi.StringArray
	// If `true`, disable the service when the resource is destroyed.
	// Defaults to `true`.
	// https://www.terraform.io/docs/providers/google/r/google_project_service.html#disable_on_destroy
	DisableOnDestroy bool
	// If `true`, services that are enabled and which depend on this service should also be disabled when this service is destroyed.
	// If `false` or unset, an error will be generated if any enabled services depend on this service when destroying it.
	// https://www.terraform.io/docs/providers/google/r/google_project_service.html#disable_dependent_services
	DisableDependentServices bool
}
