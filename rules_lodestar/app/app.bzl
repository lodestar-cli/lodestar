def app_push(name, app_config, environment, tag):
  native.genrule(
    name = name,
    srcs = [
        "//cmd/lodestar",
        app_config
    ],
    outs = ["app_push_"+name+".txt"],
    cmd_bash = "$(locations //cmd/lodestar) app push --config-path $(locations "+app_config+") --env "+environment+" --tag "+tag+" > $@",
)
