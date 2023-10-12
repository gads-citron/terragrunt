terraform {
  source = "git::ssh://git@github.com/gads-citron/i-dont-exist.git//test/fixture-download/hello-world"
}

inputs = {
  name = "terragrunt"
}
