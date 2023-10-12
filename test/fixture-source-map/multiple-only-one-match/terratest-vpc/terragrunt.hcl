terraform {
  source = "git::ssh://git@github.com/gads-citron/i-dont-exist.git//fixture-source-map/modules/vpc?ref=master"
}

inputs = {
  name = "terratest"
}
