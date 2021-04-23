# Lodestar

**Guide your applications through all of their environments**

## Overview

GitOps has allowed for development teams to control their DevOps pipelines through the tool they use alsot every day, Git.  By declaring what their code, tests, and deployments will look like through code, developers can work quickly while being reassured what reaches their end users will be in a good state.  Although this sounds nice, code doesn't just magically show up in production.  This is what lead to applications like Kubernetes, CI Tools, and GitOps Controllers (ArgoCD, Flux, etc.) being created.  By using the proper tooling, developers have been able to automate their tests and deployments through different phases (environments) before deploying into production.  A common GitOps pipeline often looks like this:

<div style="text-align:center">
<img src="https://user-images.githubusercontent.com/51719751/113178192-f6f1ae80-9213-11eb-8e34-7a94e5d87a23.png">
</div>

As stated above, there are many tools that have been made to run Continuous Integration, Continuous Deployment, and host environments for your code to deploy into.  The one area that isn't covered, and is the key difference between a GitOps pipeline and a standard pipeline, is the updating of manifests between CI and CD runs. Some teams are able to script together some commands to update manifests, while others choose to keep this as a manual step, but  both of these implementations tend to lead to potential human errors and inconsistencies between teams.

Lodestar was created to resolve this issue.  Lodestar is a tool that guides deployments to their destination in a consistent manner no matter the environment. With Lodestar, updating an environment's configuration values or promoting an application to the next environment can be boiled down to a single command.  Not only is Lodestar simple to use, it is also flexible, for it can be ran locally or within your CI/CD pipelines.  No matter where you are in your GitOps journey, Lodestar can help make sure your applications reach their environments safely.

## Implementation

### Permissions

**Note: Lodestar currently only supports authentication through personal access tokens over http.  SSH may be added later, but is currently not available.  Please see your Version Control Provider's documentation on creating personal access tokens.**

Because we want Lodestar to be able to take care of updating our manifests in our Git repository, we need to make sure that it has the correct permissions for our repositories.  To do this, you can either pass in a Git account's username and personal access token through the commandline by using the --username and --token flags, or by assigning the environment variables GIT_USER and GIT_TOKEN.  For CI pipelines, it may make sense to create a CI user in your GitHub organization and keep its credentials as a secret in your Kubernetes Cluster that can be pulled into environment variables.  With tools like Sops and Sealed Secrets, this process can be entirely maintained within your Git repository.

### AppConfigurationFile

In order for Lodestar to manage your applications, the application first needs to be registered through an appConfiguration yaml. This App Configuration is made up of three main objects: info, environmentGraph, and yamlKeys.

appConfiguration files can be created imperatively by running:

    lodestar app create

which will walk you through creating an appConfiguration.  This will automatically create and load the appConfiguration file into the lodestar home directory at ~/.lodestar/app/***appName***.yaml.  appConfigurations can also be declaratively created and placed undert the app folder in the lodestar home directory or passed through the terminal through the --config-path flag when attempting to push or promote.  It is common to include the appConfig in your code repository and pass it to lodestar when running in your pipelines.

#### Info

The Info object holds the metadata about your application so that Lodestar know's what repository holds the manifests and where the app state file is located at.  This also holds the name and the description of the app so others can easily identify what is being updated.

#### EnvironmentStateGraph

The EnvironmentStateGraph is waht helps lodestar located the files it will be updating for each environment.  It is important to keep this part up to date so that Lodestar can find the files correctly and update tehm when your run a push or a promote.

#### YamlKeys

YamlKeys is what identifies the keys for Lodestar to look for in the configuration management file that is being updated.  It is a list of strings that Lodestar will use to search with.  if you wish to replace the key tag with a new value, simply add tag to the yamlKeys list and Lodestar will now to look for it when pushing or updating.  If you wish to update a key, it must be in this list, for Lodestar will fail if it is given a Key that doesn't match the ones within this list.

**Sample appConfiguration**

    info:
      name: lodestar-folder-app-example
      type: folder
      description: this is a test app
      repoUrl: https://github.com/lodestar-cli/lodestar-folder-app-example.git
      target: main
      statePath: lodestar-folder-app-example.yaml
    environmentGraph:
    - name: dev
        srcPath: 1-dev/values.yaml
    - name: qa
      srcPath: 2-qa/values.yaml
    - name: staging
      srcPath: 3-staging/values.yaml
    - name: production
      srcPath: 4-production/values.yaml
    yamlKeys:
      - tag

**Note: Lodestar currently only allows for configuration files to be on the same branch and in the same repository.  Environments should be declared similarly to how Kustomize or Terraform handles their different environments.  Other deployment strategies like branching environments will be added at a later date**

### StateGraph
  
Although you will not have to create the stateGraph, it is important to know what it does and why it's there.  When you create your appConfig, Lodestar will go through the environments in the envGraph and create an environment in the stateGraph.  Lodestar will then go to the location in the manifest repository declared in the envGraph environment, grab the key from the manifest, and place it in the stateGraph under the proper environment.  The stateGraph is always in sync with your environments, and shouldn't be updated by hand.  The main reason behind the stateGraph is so that Lodestar can keep track of what state you want your environments to be in without needing a backend database to run.

**Sample stateGraph**

    updated: 2021-04-17T16:09:17-05:00
    environmentStateGraph:
    - name: dev
      yamlKeys: {tag: aaabbb}
    - name: qa
      yamlKeys: {tag: 12l8767}
    - name: staging
      yamlKeys: {tag: 12345}
    - name: production
      yamlKeys: {tag: 12345}
## Documentation

**Note: Currently, lodestar only supports pushing and promoting of a single application at a time.  More work may be done to expand the capabilities in the future**

### App

#### Push

NAME:
   lodestar app push - Push yaml key values to an environment

USAGE:
   In order to push keys to an environment, either a name for an App configured in ~/.lodestar
   needs to be provided with --name, or a path to an App needs to be provided with --config-path.
   Lodestar will then be able to find the App and pass the keys to the correct environment.

OPTIONS:
   --name name                                   the name of an app
   --config-path path                            the path to the app configuration file
   --environment environment, --env environment  the environment the new yaml keys will be pushed to
   --yaml-keys "key=value"                       a  comma separated "key=value" string of yaml keys to update [$YAML_KEYS]
   --output-state                                will create a local yaml file of the updated app state when set (default: false)

#### Promote

NAME:
   lodestar app promote - Promote an an environment's key values to the next environment

USAGE:
   In order to promote an environment's keys, either a name for an App configured in ~/.lodestar
   needs to be provided with --name, or a path to an App needs to be provided with --config-path.
   Lodestar will then be able to find the App and pass the keys to the correct environment.

OPTIONS:
   --name name         the name of an a
   --config-path path  the path to the a configuration file
   --src-env name      the name of the source environment
   --dest-env name     the name of the destination
   --output-state      will create a local yaml file of the updated a state when set (default: false)

#### List

NAME:
   lodestar app list - List current context Apps

USAGE:
   Will provide all the Apps within the current context as well as a description of the App.
   App names and descriptions come directly from the appInfo block in their respective App configuration file.

#### Show

NAME:
   lodestar app show - Prints information on the provided App

USAGE:
   WHen provided an App name of path to an AppConfiguration file, Show will print out both the AppConfiguration file
  as well as the current state of the Application's environments.

OPTIONS:
   --name name         the name of the a
   --config-path path  the path to the a configuration file