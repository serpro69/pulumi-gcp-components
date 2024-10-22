package gcsbackend

import "github.com/serpro69/pulumi-google-components/utils"

const (
	pkg = "gcsbackend"
)

var (
	bucket     = utils.NewResourceType(pkg, "Bucket")
	iam        = utils.NewResourceType(pkg, "Iam")
	kms        = utils.NewResourceType(pkg, "Kms")
	gcsBackend = utils.NewResourceType(pkg, "GcsBackend")
)
