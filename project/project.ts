import * as pulumi from "@pulumi/pulumi";
import * as gcp from "@pulumi/gcp";
import * as time from "@pulumiverse/time";
// tsserver: Cannot find name 'Project'.

class ProjectArgs {
  readonly billingAccount: pulumi.Input<string>;
  readonly folderId: pulumi.Input<string>;
  readonly projectId: pulumi.Input<string>;
  readonly name: pulumi.Input<string>;
  readonly autoCreateNetwork: pulumi.Input<boolean>;
  readonly labels: pulumi.Input<{ [key: string]: pulumi.Input<string> }>;
  readonly deletionPolicy: pulumi.Input<string>;
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Project extends pulumi.ComponentResource {
  public readonly main: gcp.organizations.Project;

  constructor(
    name: string,
    args: ProjectArgs,
    opts: pulumi.ComponentResourceOptions,
  ) {
    super("pgc:project:Project", name, args, opts);

    const project = new gcp.organizations.Project(
      name,
      {
        ...args,
      },
      { parent: this },
    );

    const sleep = new time.Sleep(
      name,
      { createDuration: "30s" },
      {
        parent: this,
        dependsOn: project,
        deletedWith: project,
      },
    );
    this.main = project;

    this.registerOutputs({
      project: this.main,
      wait: sleep,
    });
  }
}
