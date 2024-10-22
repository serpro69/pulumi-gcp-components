package gcsbackend

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/serpro69/pulumi-google-components/utils"
)

type GcsBackendArgs struct {
	// GCP billing account ID
	BillingAccount pulumi.StringInput
	// GCP folder ID
	FolderId pulumi.StringInput
	// GCP project ID for pulumi state
	ProjectId pulumi.StringInput
	// GCP project name for pulumi state
	ProjectName pulumi.StringInput
	// Is this a production environment?
	IsProd utils.Pair[bool, pulumi.BoolInput]

	// IAM

	// List of members to grant storage.admin role to on the project level
	IamGcsAdmins []string
	// List of members to grant storage.objectViewer role to on the project level
	IamGcsObjectViewers []string
	// List of members to grant storage.admin role to on the 'prod' folder level
	IamGcsPlStateFolderProdAdmins []string
	// List of members to grant storage.user role to on the 'prod' folder level
	IamGcsPlStateFolderProdUsers []string
	// List of members to grant storage.admin role to on the 'test' folder level
	IamGcsPlStateFolderTestAdmins []string
	// List of members to grant storage.user role to on the 'test' folder level
	IamGcsPlStateFolderTestUsers []string
}
