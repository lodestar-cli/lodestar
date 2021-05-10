load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")
load("@io_bazel_rules_docker//contrib:passwd.bzl", "passwd_entry", "passwd_file")
load("@rules_pkg//:pkg.bzl", "pkg_tar")

def update_registry_image(name, service, binary, registry, repository, base, tag):
    passwd_entry(
        username = "root",
        uid = 0,
        gid = 0,
        name = "root_"+name,
    )

    passwd_entry(
        username = "lodestar-user",
        info = "nonroot",
        uid = 1002,
        name = "user_"+name,
    )

    passwd_file(
        name = "passwd_"+name,
        entries = [
            ":root_"+name,
            ":user_"+name,
        ],
    )

    # Create a tar file containing the created passwd file
    pkg_tar(
        name = "passwd_tar_"+name,
        srcs = [":passwd_"+name],
        mode = "0o644",
        package_dir = "etc",
    )

    container_image(
        name = "image_"+ name,
        base = base,
        workdir = "/",
        entrypoint = ["/" + service],
        files = [binary],
        visibility = ["//visibility:public"],
        tars = [":passwd_tar_"+name],
        symlinks = { "/etc/passwd" : "/etc/passwd_"+name},
        user = "lodestar-user",
    )

    container_push(
        name = name,
        format = "Docker",
        image = ":image_"+ name,
        registry = registry,
        repository = repository,
        tag = tag,
    )
