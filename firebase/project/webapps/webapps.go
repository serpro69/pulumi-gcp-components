package webapps

import (
	"errors"
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/firebase"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/serpro69/pulumi-google-components/firebase/project/util"
	"github.com/serpro69/pulumi-google-components/firebase/project/vars"
)

type FirebaseProjectWebApps struct {
	pulumi.ResourceState

	Apps    firebase.WebAppMapOutput              `pulumi:"apps"`
	Domains firebase.HostingCustomDomainMapOutput `pulumi:"domains"`
}

func ConfigureWebApps(
	ctx *pulumi.Context,
	name string,
	args *vars.ProjectWebAppsArgs,
	opts ...pulumi.ResourceOption,
) (*FirebaseProjectWebApps, error) {
	// Check for mandatory arguments
	if args == nil || args.Project == nil {
		return nil, errors.New("ProjectId is mandatory to configure firebase web apps")
	}

	fb := &FirebaseProjectWebApps{}
	if err := ctx.RegisterComponentResource(util.WebApps.String(), name, fb, opts...); err != nil {
		return nil, err
	}

	fb.Apps = args.Project.ToStringOutput().ApplyT(func(projectId string) (firebase.WebAppMapOutput, error) {
		webApps := args.WebApps.ToStringArrayOutput().ApplyT(func(apps []string) (map[string]*firebase.WebApp, error) {
			am := make(map[string]*firebase.WebApp, len(apps))
			for _, app := range apps {
				var err error
				am[app], err = firebase.NewWebApp(ctx, app,
					&firebase.WebAppArgs{
						Project:     pulumi.String(projectId),
						DisplayName: pulumi.String(app),
					},
					pulumi.Parent(fb),
				)
				if err != nil {
					return nil, err
				}

				hs, err := firebase.NewHostingSite(ctx, app,
					&firebase.HostingSiteArgs{
						Project: am[app].Project,
						AppId:   am[app].AppId,
						SiteId:  pulumi.Sprintf("%s-%s", app, projectId),
					},
					pulumi.Parent(am[app]),
				)
				if err != nil {
					return nil, err
				}

				fb.Domains = args.CustomDomains.ToStringArrayMapOutput().ApplyT(func(domains map[string][]string) (map[string]*firebase.HostingCustomDomain, error) {
					dm := make(map[string]*firebase.HostingCustomDomain, len(apps))
					for _, domain := range domains[app] {
						hs.SiteId.ApplyT(func(siteId *string) error {
							dm[domain], err = firebase.NewHostingCustomDomain(ctx, fmt.Sprintf("%s$%s", app, domain),
								&firebase.HostingCustomDomainArgs{
									Project:        am[app].Project,
									SiteId:         pulumi.String(*siteId),
									CertPreference: pulumi.String("DEDICATED"),
									CustomDomain:   pulumi.String(domain),
								},
								pulumi.Parent(hs),
							)
							if err != nil {
								return err
							}
							return nil
						})
					}
					return dm, nil
				}).(firebase.HostingCustomDomainMapOutput)
			}
			return am, nil
		}).(firebase.WebAppMapOutput)
		return webApps, nil
	}).(firebase.WebAppMapOutput)

	if err := ctx.RegisterResourceOutputs(fb, pulumi.Map{
		"apps":    fb.Apps,
		"domains": fb.Domains,
	}); err != nil {
		return nil, err
	}

	return fb, nil
}
