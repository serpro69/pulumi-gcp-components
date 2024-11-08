package project

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/serpro69/pulumi-google-components/firebase/project/util"
	"github.com/serpro69/pulumi-google-components/firebase/project/vars"
)

// FirebaseProjectIAM is a struct that represents IAM of a given firebase project
type FirebaseProjectIam struct {
	pulumi.ResourceState
}

func configureIAM(
	ctx *pulumi.Context,
	name string,
	projectId string,
	projectNumber string,
	args *vars.ProjectIamArgs,
	opts ...pulumi.ResourceOption,
) (*FirebaseProjectIam, error) {
	fpIam := &FirebaseProjectIam{}
	if err := ctx.RegisterComponentResource(util.Iam.String(), name, fpIam, opts...); err != nil {
		return nil, err
	}

	dsa, err := compute.GetDefaultServiceAccount(ctx,
		&compute.GetDefaultServiceAccountArgs{
			Project: pulumi.StringRef(projectId),
		},
		pulumi.Parent(fpIam),
	)
	if err != nil {
		return nil, err
	}

	csa := args.ComputeServiceAccountRoles.ToStringArrayOutput().ApplyT(func(roles []string) ([]*projects.IAMMember, error) {
		var mm []*projects.IAMMember
		for _, role := range roles {
			m, err := addMemberRole(ctx, fpIam, projectId, dsa.Member, role)
			if err != nil {
				return nil, err
			}
			mm = append(mm, m)
		}
		return mm, nil
	})

	psa := args.PubSubServiceAccountRoles.ToStringArrayOutput().ApplyT(func(roles []string) ([]*projects.IAMMember, error) {
		sa := fmt.Sprintf("serviceAccount:service-%s@gcp-sa-pubsub.iam.gserviceaccount.com", projectNumber)
		var mm []*projects.IAMMember
		for _, role := range roles {
			m, err := addMemberRole(ctx, fpIam, projectId, sa, role)
			if err != nil {
				return nil, err
			}
			mm = append(mm, m)
		}
		return mm, nil
	})

	admins := args.FirebaseAdminMembers.ToStringArrayOutput().ApplyT(func(members []string) ([]*projects.IAMMember, error) {
		var mm []*projects.IAMMember
		for _, member := range members {
			m, err := addMemberRole(ctx, fpIam, projectId, member, "firebase.admin")
			if err != nil {
				return nil, err
			}
			mm = append(mm, m)
		}
		return mm, nil
	})

	viewers := args.FirebaseViewerMembers.ToStringArrayOutput().ApplyT(func(members []string) ([]*projects.IAMMember, error) {
		var mm []*projects.IAMMember
		for _, member := range members {
			m, err := addMemberRole(ctx, fpIam, projectId, member, "firebase.viewer")
			if err != nil {
				return nil, err
			}
			mm = append(mm, m)
		}
		return mm, nil
	})

	// // wait for iammember outputs before sleeping
	// w := pulumi.All(csaOut, psaOut).ApplyT(func(ss []interface{}) (*time.Sleep, error) {
	// 	// wait for members to be added
	// 	wfs, err := time.NewSleep(ctx, fmt.Sprintf("%s/wait", name),
	// 		&time.SleepArgs{
	// 			CreateDuration: pulumi.String("30s"),
	// 			Triggers: pulumi.StringMap{
	// 				// panic: applier's first input parameter must be assignable from []*projects.IAMMember, got []string
	// 				//     applier defined at /home/sergio/Projects/personal/pulumi-google-components/firebase/project/iam.go:97
	// 				"members": fpIam.IAMMemberArray.ToIAMMemberArrayOutput().ApplyT(func(mm []*projects.IAMMember) (string, error) {
	// 					members := make([]string, len(mm))
	// 					for i, m := range mm {
	// 						m.Member.ToStringOutput().ApplyT(func(s string) error {
	// 							members[i] = s
	// 							return nil
	// 						})
	// 					}
	// 					return strings.Join(members, ","), nil
	// 				}).(pulumi.StringOutput),
	// 			},
	// 		},
	// 		pulumi.Parent(fpIam),
	// 		pulumi.DeletedWith(fpIam),
	// 		pulumi.DependsOn(press),
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return wfs, nil
	// }).(time.SleepOutput)

	if err := ctx.RegisterResourceOutputs(fpIam, pulumi.Map{
		"defaultComputeSA":           pulumi.String(dsa.Member),
		"computeServiceAccountRoles": csa,
		"pubSubServiceAccountRoles":  psa,
		"adminMembers":               admins,
		"viewerMembers":              viewers,
		// "wait":             w,
		// "triggers": w.ApplyT(func(sleep *time.Sleep) (pulumi.StringPtrOutput, error) {
		// 	so := sleep.Triggers.ApplyT(func(triggers map[string]string) (*string, error) {
		// 		s := triggers["members"]
		// 		return &s, nil
		// 	}).(pulumi.StringPtrOutput)
		// 	return so, nil
		// }).(pulumi.StringPtrOutput),
	}); err != nil {
		return nil, err
	}

	return fpIam, nil
}

func addMemberRole(ctx *pulumi.Context, parent pulumi.Resource, projectId, member, role string) (*projects.IAMMember, error) {
	m, err := projects.NewIAMMember(ctx, fmt.Sprintf("%v/%v", member, role),
		&projects.IAMMemberArgs{
			Project: pulumi.String(projectId),
			Role:    pulumi.String(fmt.Sprintf("roles/%v", role)),
			Member:  pulumi.String(member),
		},
		pulumi.Parent(parent),
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}
