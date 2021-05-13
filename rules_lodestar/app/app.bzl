def _app_push_impl(ctx):
    ctx.actions.run_shell(
        inputs = [ctx.file.app_config, ctx.file.yaml_keys, ctx.file.token, ctx.file.username, ctx.file._lodestar, ctx.file.wdir],
        outputs = [ctx.outputs.execute],
        execution_requirements = {
            "no-sandbox":"1",
            "no-remote":"1"
        },
        command = """cat {wdir} | tr '\n' '/' > {out} && echo -n {lodestar} app push --env {environment}  >> {out} && 
         echo -n \" --config-path \" >> {out} && cat {wdir} | tr '\n' '/' >> {out} && echo -n {ac} >> {out} &&
         echo -n \" --username \" >> {out} && cat {username} | tr '\n' ' ' >> {out} &&
         echo -n \"--token \" >> {out} && cat {token} | tr '\n' ' ' >> {out} &&
         echo -n \"--yaml-keys \\\"\" >> {out} && cat {yk} | tr -d '\n' >> {out} && echo -n \"\\\"\" >> {out}
         """.format(
            wdir= ctx.file.wdir.path,
            lodestar = ctx.file._lodestar.path,
            environment = ctx.attr.environment,
            yk = ctx.file.yaml_keys.path,
            ac = ctx.file.app_config.path,
            username = ctx.file.username.path,
            token = ctx.file.token.path,
            out = ctx.outputs.execute.path
        ),
    )
    return [DefaultInfo(executable=ctx.outputs.execute)]


_app_push = rule(
    implementation = _app_push_impl,
    executable = True,
    attrs = {
        "_lodestar": attr.label(
            allow_single_file = True,
            executable = True,
            cfg = "exec",
            default = Label("@lodestar//cmd/lodestar:lodestar"),
        ),
        "out":  attr.output(doc='', mandatory=False),
        "app_config": attr.label(allow_single_file = True),
        "environment": attr.string(mandatory=True),
        "yaml_keys": attr.label(allow_single_file = True),
        "username": attr.label(allow_single_file = True),
        "token": attr.label(allow_single_file = True,),
        "wdir": attr.label(allow_single_file = True)
    },
    outputs = {
        "execute": "execute_lodestar.sh",
    }
)

def _get_stamp(name, stamp):
    if "{" in stamp:
        s = stamp[1:len(stamp)-1]
        native.genrule(
            name = name,
            stamp = 1,
            outs = [name+".txt"],
            cmd = "cat bazel-out/stable-status.txt | grep "+s+" | awk '/"+s+"/{print $$NF}' > $@",
        )
    else:
        native.genrule(
            name = name,
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

    _app_push(
        name= name,
        app_config = app_config,
        environment = environment,
        yaml_keys = name+"_yk",
        username = name+"_username",
        token = name+"_token",
        wdir = name+"_workspace"
    )