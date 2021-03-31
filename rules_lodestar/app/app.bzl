def app_push(name, lodestar, app_config, environment, tag):
  native.genrule(
    name = name,
    srcs = [
        lodestar,
        app_config
    ],
    outs = ["app_push_"+name+".txt"],
    cmd_bash = "$(locations "+lodestar+") app push --config-path $(locations "+app_config+") --env "+environment+" --tag "+tag+" > $@",
)

def app_promote(name, lodestar, app_config, src_env, dest_env):
  native.genrule(
    name = name,
    srcs = [
        lodestar,
        app_config
    ],
    outs = ["app_push_"+name+".txt"],
    cmd_bash = "$(locations "+lodestar+") app promote --config-path $(locations "+app_config+") --src-env "+src_env+" --dest-env "+dest_env+" > $@",
)