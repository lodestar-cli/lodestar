def app_push(name, lodestar, app_config, environment, tag):
  native.genrule(
    name = name,
    srcs = [
        lodestar,
        app_config
    ],
    local=True,
    outs = ["app_push_"+name+".txt"],
    cmd_bash = "$(locations "+lodestar+") app push --config-path $(locations "+app_config+") --env "+environment+" --tag "+tag+" > $@",
)
