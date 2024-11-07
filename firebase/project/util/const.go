package util

import "github.com/serpro69/pulumi-google-components/utils"

const pkg = "firebase/project"

var (
	Project = utils.NewResourceType(pkg, "Project")
	Iam     = utils.NewResourceType(pkg, "Iam")
)
