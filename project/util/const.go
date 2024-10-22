package util

import "github.com/serpro69/pulumi-google-components/utils"

const pkg = "project"

var (
	Project  = utils.NewResourceType(pkg, "Project")
	Services = utils.NewResourceType(pkg, "Services")
)
