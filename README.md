# Lodestar

**Guide your applications through all of their environments**

## Overview

Based on [Continuous Delivery](https://continuousdelivery.com/),  GitOps culture has allowed for development teams to truly connect their audiences to the code they write every day.  By declaring what their code, tests, and deployments will look like through Git, developers can work quickly with the reassurance that what reaches production will always be in a good state.  Although this sounds nice, code doesn't just magically show up in production.  This is what lead to tools like Kubernetes, CI Servers, and GitOps Controllers (ArgoCD, Flux, etc.) being created.  By using the proper tooling, developers have been able to automate their tests and deployments through different phases (environments) before deploying into production.  A common GitOps pipeline often looks like this:

<div style="text-align:center">
<img src="https://user-images.githubusercontent.com/51719751/113178192-f6f1ae80-9213-11eb-8e34-7a94e5d87a23.png">
</div>

As stated above, there are many tools that have been made to run Continuous Integration, Continuous Deployment, and host environments for your code to deploy into.  The one area that isn't covered, and is the key difference between a GitOps and a standard pipeline, is the updating of manifests between CI and CD runs.  There are different strategies for implementing this.  Some teams are able to script together some commands to update manifests, while others choose to keep this as a manual step.  Both these implementations lead to potential human errors and inconsistencies between teams.

In comes Lodestar.  Lodestar is a tool that guides deployments to their destination in a consistent manner no matter the environment. With Lodestar, pushing a new image tag to an environment or promoting an image to the next environment can be boiled down to a single command.  Lodestar can be as simple or complex as you need it to be, and can grow with your deployment flow.  Not only does Lodestar scale with your pipeline, it is flexible to your deployment strategy, for it can be ran locally or within your automation pipeline.

## Implementation

### Permissions

**Note: Lodestar currently only supports authentication through personal access tokens over http.  SSH may be added later, but is currently not available.  Please see your Version Control Provider's documentation on creating personal access tokens.**

Because we want Lodestar to be able to take care of updating our manifests in our Git repository, we need to make sure that it has the correct permissions for our repositories.  To do this, you can either pass in a Git account's username and personal access token through the commandline by using the --username and --token flags, or by placing them under the environment variables GIT_USER and GIT_TOKEN.  For CI pipelines, it may make sense to create a CI user and keep its credentials as a secret in your Kubernetes Cluster that can be pulled into an environment variable.  Secure secret implementation can also be done through Sealed Secrets or Sops.

### Tagging

The main goal of Lodestar is to be useful without hendering the readability of your deployments.  Both Helm and Kustomize have a learning curve to understand what is going on in their manifests, and Lodestar doesn't want to make that any harder to read.  The only requirement that Lodestar has is that you need to wrap whatever is holding your tag with ###lodestar###.  An example would look like:

    ###lodestar###
    tag: v1.2.3
    ###lodestar###

This takes advantage of yaml's ability to use comments and won't impact how you deploy your manifests.


### AppConfiguration

In order for Lodestar to manage your application's images, the application first needs to be registered through an appConfiguration yaml.  This App Configuration is made up of two main objects: an appInfo object and an envGraph.  The appInfo object holds the metadata for the app.  This includes the name of the app, the url the app's manifests are located at, the path to the stateGraph yaml, and other common information.  The envGraph holds the information needed for Lodestar to correctly update your app's environments.  This tells Lodestar how many environments are in the app, and where to locate the manifests to update the tags with.

**Sample appConfiguration**

    appInfo:
      name: lodestar-folder-app-example
      type: folder
      description: this is a test app
      repoUrl: https://github.com/lodestar-cli/lodestar-folder-app-example.git
      target: main
      statePath: lodestar-folder-app-example.yaml
    envGraph:
      - name: dev
      srcPath: 1-dev/values.yaml
      - name: qa
      srcPath: 2-qa/values.yaml
      - name: staging
      srcPath: 3-staging/values.yaml
      - name: production
      srcPath: 4-production/values.yaml

**Note: Lodestar currently only allows for configuration files to be on the same branch and in the same repository.  Environments should be declared similarly to how Kustomize or Terraform handles their different environments.  Other deployment strategies like branching environments will be added at a later date**

appConfiguration files can be created imperatively by running:

    lodestar app create

which will walk you through creating an appConfiguration.  This will automatically create and load the appConfiguration file into the lodestar home directory at ~/.lodestar/app/***appName***.yaml.  appConfigurations can also be declaratively created and placed undert the app folder in the lodestar home directory or passed through the terminal through the --config-path flag when attempting to push or promote.  It is common to include the appConfig in your code repository and pass it to lodestar when running in your pipelines.

### StateGraph
  
Although you will not have to create the stateGraph, it is important to know what it does and why it's there.  When you create your appConfig, Lodestar will go through the environments in the envGraph and create an environment in the stateGraph.  Light house will then go to the location in the manifest repository declared in the envGraph environment, grab the tag from the manifest, and place it in the stateGraph under the proper environment.  The stateGraph is always in sync with your environments, and shouldn't be updated by hand.  The main reason behind the stateGraph is so that Lodestar can keep track of what state you want your environments to be in without needing a backend database to run.

**Sample stateGraph**

    state:
      - environment: dev
      tag: v1.1.1
      - environment: qa
      tag: v1.1.0
      - environment: staging
      tag: v1.1.0
      - environment: production
      tag: v1.1.0                      

## Documentation

**Note: Currently, lodestar only supports pushing and promoting of a single application at a time.  More work may be done to expand the capabilities in the future**

### App

#### Push

Pushing allows for you to define a new tag for an environment.  This is handy for when you just created a new image and want it to go into a development environment for integration testing.  Pushing is typically done after a new image is created and the tag isn't in an environment yet.  Either a name for an app currently loaded in the home directory or a path to a configuration file must be provided.

    lodestar app push [options]

    Options:
    --name name                                   the name of an app
    --config-path path                            the path to the app configuration file
    --environment environment, --env environment  the environment the tag will be pushed to
    --tag tag                                     the tag for the new image
    --output-state                                will create a local yaml file of the updated app state when set

#### Promote

Promoting an image between environments allows for the promotion of an image without having to know the tag.  This will grab the tag in a source environment and promote it to a destination environment.  Either a name for an app currently loaded in the home directory or a path to a configuration file must be provided.

    lodestar app promote [options]

    Options:
    --name name                                   the name of an app
    --config-path path                            the path to the app configuration file
    --environment environment, --env environment  the environment the tag will be pushed to
    --tag tag                                     the tag for the new image
    --output-state                                will create a local yaml file of the updated app state when set

#### List

Will list the name and description of all apps currently in the ~/.lodestar/app directory

    lodestar app list

#### Show

Will show the appConfiguration of an app when provided a name.

    lodestar app show --name [name]