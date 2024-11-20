package project

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/firebase"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/serpro69/pulumi-google-components/firebase/project/util"
	"github.com/serpro69/pulumi-google-components/firebase/project/vars"
)

type FirebaseProjectWebApps struct {
	pulumi.ResourceState
}

func configureWebApps(
	ctx *pulumi.Context,
	name string,
	projectId string,
	args *vars.ProjectWebAppsArgs,
	opts ...pulumi.ResourceOption,
) (*FirebaseProjectWebApps, error) {
	fbWebApps := &FirebaseProjectWebApps{}
	if err := ctx.RegisterComponentResource(util.WebApps.String(), name, fbWebApps, opts...); err != nil {
		return nil, err
	}

	var aac []*firebase.GetWebAppConfigResult
	var acd []*firebase.HostingCustomDomain
	webApps := args.WebApps.ToStringArrayOutput().ApplyT(func(apps []string) ([]*firebase.WebApp, error) {
		var aa []*firebase.WebApp
		for _, app := range apps {
			a, err := firebase.NewWebApp(ctx, app,
				&firebase.WebAppArgs{
					Project:     pulumi.String(projectId),
					DisplayName: pulumi.String(app),
				},
				pulumi.Parent(fbWebApps),
			)
			if err != nil {
				return nil, err
			}
			aa = append(aa, a)

			a.AppId.ApplyT(func(appId string) error {
				ac, err := firebase.GetWebAppConfig(ctx,
					&firebase.GetWebAppConfigArgs{
						Project:  pulumi.StringRef(projectId),
						WebAppId: appId,
					},
					pulumi.Parent(a),
				)
				if err != nil {
					return err
				}
				aac = append(aac, ac)
				return nil
			})

			hs, err := firebase.NewHostingSite(ctx, app,
				&firebase.HostingSiteArgs{
					Project: a.Project,
					AppId:   a.AppId,
					SiteId:  pulumi.Sprintf("%s-%s", app, projectId),
				},
				pulumi.Parent(a),
			)
			if err != nil {
				return nil, err
			}

			args.CustomDomains.ToStringArrayMapOutput().ApplyT(func(domains map[string][]string) error {
				for _, domain := range domains[app] {
					hs.SiteId.ApplyT(func(siteId *string) error {
						d, err := firebase.NewHostingCustomDomain(ctx, fmt.Sprintf("%s$%s", app, domain),
							&firebase.HostingCustomDomainArgs{
								Project:        a.Project,
								SiteId:         pulumi.String(*siteId),
								CertPreference: pulumi.String("DEDICATED"),
								CustomDomain:   pulumi.String(domain),
							},
							pulumi.Parent(hs),
						)
						if err != nil {
							return err
						}
						acd = append(acd, d)
						return nil
					})
				}
				return nil
			})
		}
		return aa, nil
	})

	if err := ctx.RegisterResourceOutputs(fbWebApps, pulumi.Map{
		"webApps": webApps,
	}); err != nil {
		return nil, err
	}

	return fbWebApps, nil
}
