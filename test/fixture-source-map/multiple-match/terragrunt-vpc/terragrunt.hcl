terraform {
  source = "git::ssh://git@github.com/gads-citron/another-dont-exist.git//fixture-source-map/modules/vpc?ref=master"
}

inputs = {
  name = "terragrunt"
}
