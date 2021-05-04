# Lodestar Bazel Rules

This folder provides some Bazel genrules wrapped in macros to run in your bazel builds.  Bazel is often used to a,build, and deploy a binary, and these can be added to the end of that flow to keep everything ran with a single tool.  This also allows for Lodestar to be ran without having it installed on a local machine.  As long as Bazel is installed and Lodestar is declared in the WORKSPACE, then it should be able to run the same as if it were called from the commandline.

### WORKSPACE

To run Lodestar in your BUILD file, you must fist checkout the repository in your workspace:

    load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

    git_repository(
        name = "lodestar",
        commit = "506d99e4b7e9efc38969f1e9e3cfb2bde9e9d791",
        remote = "https://github.com/lodestar-cli/lodestar.git",
    )

### BUILD.bazel

in order to push or promote with Bazel, you need to load the app.bzl file into your BUILD file.  Once loaded, you then can run either app_push or app_promote.

    load(":app.bzl", "app_push", "app_promote")

    app_push(
        name = "push_example",
        app_config = "lodestar-folder-app-example.yaml",
        environment = "dev",
        lodestar = "//cmd/lodestar",
        tag = "$(tag)",
    )

    app_promote(
        name = "promote_example",
        lodestar = "//cmd/lodestar",
        app_config = "lodestar-folder-app-example.yaml",
        src_env = "dev",
        dest_env = "qa",
    )