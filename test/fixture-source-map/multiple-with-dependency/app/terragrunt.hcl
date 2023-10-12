terraform {
  source = "git::ssh://git@github.com/gads-citron/another-dont-exist.git//fixture-source-map/modules/app?ref=master"
}

dependency "vpc" {
  config_path = "../vpc"
}

inputs = {
  name   = "terragrunt"
  vpc_id = dependency.vpc.outputs.vpc_id
}
