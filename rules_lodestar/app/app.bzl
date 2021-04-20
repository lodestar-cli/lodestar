def app_push(name, app_config, environment, yaml_keys):
  native.genrule(
    name = name,
    executable = True,
    srcs = [
        app_config
    ],
    exec_tools = ["//cmd/lodestar:lodestar"],
    outs = ["app_push_"+name+".sh"],
    cmd_bash = "$(location //cmd/lodestar:lodestar) app push --config-path $(locations "+app_config+") --env "+environment+" --yaml-keys "+yaml_keys+" && echo \"echo Lodestar Push Complete!\" > $@",
)