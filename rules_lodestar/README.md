# Lodestar Bazel Rules

This folder provides some Bazel genrules wrapped in macros to run in your bazel builds.  Bazel is often used to a,build, and deploy a binary, and these can be added to the end of that flow to keep everything ran with a single tool.  This also allows for Lodestar to be ran without having it installed on a local machine.  As long as Bazel is installed and Lodestar is declared in the WORKSPACE, then it should be able to run the same as if it were called from the commandline.

**NOTE: Do not run these rules on Windows. They will build appropriately, however, they will throw errors if attempting to bazel run them.  Please use either Linux or MacOS to run these rules.**

## WORKSPACE

To run Lodestar in your BUILD file, you must fist checkout the repository in your workspace:

    load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

    git_repository(
        name = "lodestar",
        commit = "506d99e4b7e9efc38969f1e9e3cfb2bde9e9d791",
        remote = "https://github.com/lodestar-cli/lodestar.git",
    )

## Stamping
The lodestar bazel rules use stamping to access secret information like your token as well as info that may change such as yaml-keys that need to be pushed.  To do this, you will need to make sure two files exist in your repository at the WORKSPACE level before attempting to use any of these rules.

### .bazelrc
The .bazelrc file allows you to add add commands to your cli without having to retype them every time.  Specifically, we need to make sure when you run bazel build or run, that stamps are properly being applied to  your commands.  If the .bazelrc file doesn't already exist, create it at the WORKSPACE level and add

    build --stamp --workspace_status_command=./tools/bazel_stamp_vars.sh
    run --workspace_status_command=./tools/bazel_stamp_vars.sh

The important detail is that you added stamp and workspace_status_command to your build and run commands.  The sh file doesn't need to be named that can could potentially be a different type of script if need be.  For ease of use though, you should use the .sh file that is provided in this readme.

### bazel_stamp_vars.sh
Here we are providing the commands that are adding our variables at the time you run your bazel commands.  This file adds information to the stable-status and volatile-status files in your bazel cache.  In order to use the lodestar bazel rules, your .sh file should have the following

    # /bin/bash
    cat << EOF
    STABLE_GIT_USER $(echo $GIT_USER) //picks up your username that's set as an env var
    STABLE_GIT_TOKEN $(echo $GIT_TOKEN) //picks up your token that's set as an env var
    STABLE_YAML_KEYS $(echo $YAML_KEYS) //picks up your yaml-keys that's set as an env var
    STABLE_WORKSPACE_DIR $(pwd) //gets the workspace dir
    EOF

Without these commands, your lodestar build and runs will be unsucessful.  Although you don't need to follow these exact names, you must always keep the STABLE_ section, for the lodestar rules look into the stable-volatile file for information.

## BUILD.bazel

in order to push with Bazel, you need to load the app.bzl file into your BUILD file.  Once loaded, you then can then configure your app push.

    load(":app.bzl", "app_push")

    app_push(
        name = "push",
        app_config = "lodestar-folder-app-example.yaml",
        environment = "dev",
        token = "{STABLE_GIT_TOKEN}",
        username = "{STABLE_GIT_USER}",
        yaml_keys = "{STABLE_YAML_KEYS}",
    )

## Executing

If you run:

    bazel build //:<path-to-target>

Bazel will pass, but you will only be building the command that will execute your lodestar command.  To successfully run lodestar, you will need to do bazel run"

    bazel run //:<path-to-target>