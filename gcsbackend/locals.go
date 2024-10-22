package gcsbackend

import (
	"fmt"

	random "github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type locals struct {
	kmsLocation    pulumi.StringInput
	kmsKeyRingName pulumi.StringInput
	kmsKeyName     pulumi.StringInput
	statePrefix    pulumi.StringInput
}

func newLocals(ctx *pulumi.Context, opts ...pulumi.ResourceOption) (*locals, error) {
	l := &locals{
		kmsLocation:    pulumi.String("eur4"),
		kmsKeyRingName: pulumi.String("pulumi-state"),
	}

	if p, err := statePrefix(ctx, "main", opts...); err != nil {
		return nil, err
	} else {
		l.statePrefix = p
		l.kmsKeyName = p
	}

	return l, nil
}

/*
Returns id and name string inputs for a GCP project, based on prefix and postfix arguments.

IF project !isProd, THEN name will have the same value as id.

Returns an error IF len(prefix) < 6
*/
func statePrefix(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (pulumi.StringInput, error) {
	r, err := random.NewRandomId(ctx, name,
		&random.RandomIdArgs{
			ByteLength: pulumi.Int(8),
		},
		opts...,
	)
	if err != nil {
		return nil, err
	}
	// if err != nil {
	// 	if r, e := random.GetRandomId(ctx, "prefix", r.ID(), nil); e != nil {
	// 		return nil, e
	// 	} else if r != nil {
	// 		return r.Hex.ApplyT(func(s string) string { return fmt.Sprintf("%s-tfstate", s) }).(pulumi.StringInput), nil
	// 	} else {
	// 		return nil, err
	// 	}
	// }
	id := r.Hex.ApplyT(func(s string) string { return fmt.Sprintf("%s-plstate", s) }).(pulumi.StringInput)
	return id, nil
}
