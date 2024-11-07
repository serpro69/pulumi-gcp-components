package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-time/sdk/go/time"
	"github.com/serpro69/pulumi-google-components/project/util"
	"github.com/serpro69/pulumi-google-components/project/vars"
)

// ProjectServices is a struct that represents a GCP project with an array of enabled services
type ProjectServices struct {
	pulumi.ResourceState

	projects.ServiceArray
}

func ActivateApis(
	ctx *pulumi.Context,
	name string,
	args *vars.ProjectServicesArgs,
	opts ...pulumi.ResourceOption,
) (*ProjectServices, error) {
	// Check for mandatory arguments
	if args == nil || args.ProjectId == nil {
		return nil, errors.New("ProjectId is mandatory")
	}

	ps := &ProjectServices{}
	err := ctx.RegisterComponentResource(util.Services.String(), name, ps, opts...)
	if err != nil {
		return nil, err
	}

	// TODO: wip
	// enable APIs
	// pouts := pulumi.All(args.ActivateApis).ApplyT(func(apis []string) ([]string, error) {
	// 	for _, api := range apis {
	// 		s, err := projects.NewService(ctx, api,
	// 			&projects.ServiceArgs{
	// 				Project:                  args.ProjectId,
	// 				Service:                  pulumi.String(api),
	// 				DisableOnDestroy:         args.DisableOnDestroy,
	// 				DisableDependentServices: args.DisableDependentServices,
	// 			},
	// 			pulumi.Parent(ps))
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		ps.Services = append(ps.Services, s)
	// 		press = append(press, s)
	// 	}
	// 	return apis, nil
	// }).(pulumi.StringArrayOutput)

	// TODO: another wip with ActivateApis as []string
	// don't like the double loop though
	// var press []pulumi.Resource
	// pouts := pulumi.All(args.ActivateApis).ApplyT(func(apis []interface{}) ([]projects.ServiceOutput, error) {
	// 	out := make([]projects.ServiceOutput, 0)
	// 	for _, api := range apis {
	// 		aa := api.([]string)
	// 		for _, a := range aa {
	// 			s, err := projects.NewService(ctx, a,
	// 				&projects.ServiceArgs{
	// 					Project:                  args.ProjectId,
	// 					Service:                  pulumi.String(a),
	// 					DisableOnDestroy:         pulumi.Bool(args.DisableOnDestroy),
	// 					DisableDependentServices: pulumi.Bool(args.DisableDependentServices),
	// 				},
	// 				pulumi.Parent(ps))
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			ps.Services = append(ps.Services, s)
	// 			press = append(press, s)
	// 			out = append(out, s.ToServiceOutput())
	// 		}
	// 	}
	// 	return out, nil
	// })

	var press []pulumi.Resource

	// enable APIs
	pouts := args.ActivateApis.ToStringArrayOutput().ApplyT(func(apis []string) ([]projects.ServiceOutput, error) {
		var oo []projects.ServiceOutput
		for _, a := range apis {
			s, err := projects.NewService(ctx, a,
				&projects.ServiceArgs{
					Project:                  args.ProjectId,
					Service:                  pulumi.String(a),
					DisableOnDestroy:         pulumi.Bool(args.DisableServicesOnDestroy),
					DisableDependentServices: pulumi.Bool(args.DisableDependentServices),
				},
				pulumi.Parent(ps),
				pulumi.DeletedWith(ps),
			)
			if err != nil {
				return nil, err
			}
			ps.ServiceArray = append(ps.ServiceArray, s)
			press = append(press, s)
			oo = append(oo, s.ToServiceOutput())
		}
		return oo, nil
	})

	// wait for services outputs before sleeping
	// credits: https://www.pulumi.com/ai/conversations/0225f449-28f4-4d5d-bbd6-e05673d76a86
	wfs := pulumi.All(pouts).ApplyT(func(ss []interface{}) (*time.Sleep, error) {
		// wait for services to be enabled
		wfs, err := time.NewSleep(ctx, fmt.Sprintf("%s/wait-for-services", name),
			&time.SleepArgs{
				CreateDuration: pulumi.String("30s"),
				Triggers: pulumi.StringMap{
					"services": args.ActivateApis.ToStringArrayOutput().ApplyT(func(apis []string) (string, error) {
						return strings.Join(apis, ","), nil
					}).(pulumi.StringOutput),
				},
				// TODO: for args.ActivateApis as []string
				// Triggers: pulumi.StringMap{
				// 	"services": pulumi.String(strings.Join(args.ActivateApis, ",")),
				// },
			},
			pulumi.Parent(ps),
			pulumi.DeletedWith(ps),
			pulumi.DependsOn(press),
		)
		if err != nil {
			return nil, err
		}
		return wfs, nil
	}).(time.SleepOutput)

	// always register outputs, even if they're not used
	// https://www.pulumi.com/docs/iac/concepts/resources/components/#registering-component-outputs
	if err := ctx.RegisterResourceOutputs(ps, pulumi.Map{
		"projectId": args.ProjectId.ToStringOutput(),
		"services":  ps.ServiceArray,
		"waits":     wfs,
		"triggers": wfs.ApplyT(func(sleep *time.Sleep) (pulumi.StringPtrOutput, error) {
			so := sleep.Triggers.ApplyT(func(triggers map[string]string) (*string, error) {
				s := triggers["services"]
				return &s, nil
			}).(pulumi.StringPtrOutput)
			return so, nil
		}).(pulumi.StringPtrOutput),
		// "waits": wfs.Triggers.ApplyT(func(triggers map[string]string) (*string, error) {
		// 	s := triggers["services"]
		// 	return &s, nil
		// }).(pulumi.StringPtrOutput),
	}); err != nil {
		return nil, err
	}

	return ps, nil
}
