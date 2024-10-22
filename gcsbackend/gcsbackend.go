package gcsbackend

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type GcsBackend struct {
	pulumi.ResourceState
	name   string
	Bucket *Bucket
}

func NewGcsBackend(
	ctx *pulumi.Context,
	name string,
	args *GcsBackendArgs,
	opts ...pulumi.ResourceOption,
) (*GcsBackend, error) {
	gcsb := &GcsBackend{name: name}

	if err := ctx.RegisterComponentResource(gcsBackend.String(), name, gcsb, opts...); err != nil {
		return nil, err
	}

	l, err := newLocals(ctx, pulumi.Parent(gcsb))
	if err != nil {
		return nil, err
	}

	k, err := setupKms(ctx, name, args, l, pulumi.Parent(gcsb))
	if err != nil {
		return nil, err
	}

	iam, err := setupIam(ctx, name, args, pulumi.Parent(gcsb))
	if err != nil {
		return nil, err
	}

	if b, err := newGcsBucket(ctx, name, args, l, k,
		pulumi.Parent(gcsb),
		pulumi.DependsOn([]pulumi.Resource{iam}), // must handle GCS SA permissions first to use CMK for encryption
	); err != nil {
		return nil, err
	} else {
		gcsb.Bucket = b
	}

	return gcsb, nil
}
