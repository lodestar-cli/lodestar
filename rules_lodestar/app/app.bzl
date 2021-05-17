load("@io_bazel_rules_go//go:def.bzl", "go_binary")

def _get_stamp(name, stamp):
    if "{" in stamp:
        s = stamp[1:len(stamp)-1]
        native.genrule(
            name = name,
            stamp = 1,
            local = True,
            outs = [name+".txt"],
            cmd_ps = " (cat bazel-out\\stable-status.txt | Select-String  "+s+" | Out-String).Trim() > $@",
            cmd = "cat bazel-out/stable-status.txt | grep "+s+" > $@",
        )
    else:
        native.genrule(
            name = name,
            local = True,
            outs = [name+".txt"],
            cmd = "echo "+stamp+" > $@",
        )

def app_push(name, yaml_keys, token, username, app_config, environment):

    _get_stamp(
        name=name+"_yk",
        stamp=yaml_keys
    )

    _get_stamp(
        name=name+"_token",
        stamp=token
    )

    _get_stamp(
        name=name+"_username",
        stamp=username
    )

    _get_stamp(
        name=name+"_workspace",
        stamp="{WORKSPACE_DIR}"
    )

    native.genrule(
        name = name+"_create_input",
        srcs =[name+"_yk", name+"_token", name+"_username", name+"_workspace"],
        outs = ["stamps.txt"],
        cmd = """
for SRC in $(SRCS)
do
	cat $$SRC >> $@
done
""",
    )


    native.genrule(
        name = name+"_lodestar",
        srcs = [app_config, "@lodestar//cmd/lodestar:lodestar"],
        outs = ["lodestar.txt"],
        cmd = """echo \"STABLE_LODESTAR_DIR $(location @lodestar//cmd/lodestar:lodestar)\" >> $@ &&
         echo \"STABLE_APPCONFIG_DIR $(location """+app_config+""")\" >> $@ &&
         echo -n \"STABLE_ENVIRONMENT """+environment+"""\" >> $@""",
    )


    go_binary(
        name = name,
        srcs = ["@lodestar//rules_lodestar/app:push.go"],
        data = ["stamps.txt", "lodestar.txt"]
    )




