package vars

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ProjectWebAppsArgs struct {
	// Project ID to enable APIs on.
	// Mandatory value. An error will be returned if ProjectId is not set.
	ProjectId pulumi.StringInput
	// The list of web apps to create within the project
	WebApps pulumi.StringArray
	// A map of WebApps associated with a list of custom domains
	CustomDomains pulumi.StringArrayMap
}
