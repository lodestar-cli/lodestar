def app_push(name, lodestar, app_config, environment, yaml_keys):
  native.genrule(
    name = name,
    executable = True,
    srcs = [
        app_config
    ],
    exec_tools = [lodestar],
    outs = ["app_push_"+name+".sh"],
    cmd_bash = "$(location "+lodestar+") app push --config-path $(locations "+app_config+") --env "+environment+" --yaml-keys "+yaml_keys+" && echo \"echo Lodestar Push Complete!\" > $@",
)

def app_promote(name, lodestar, app_config, src_env, dest_env):
  native.genrule(
    name = name,
    executable = True,
    srcs = [
        app_config
    ],
    exec_tools = [lodestar],
    outs = ["app_promote_"+name+".sh"],
    cmd_bash = "$(locations "+lodestar+") app promote --config-path $(locations "+app_config+") --src-env "+src_env+" --dest-env "+dest_env+" && echo \"echo Lodestar Promote Complete!\" > $@",
)