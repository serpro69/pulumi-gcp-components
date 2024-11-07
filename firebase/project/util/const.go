package util

import "github.com/serpro69/pulumi-google-components/utils"

const pkg = "firebase"

var (
	Project = utils.NewResourceType(pkg, "Project")
	Iam     = utils.NewResourceType(pkg, "Iam")
)
