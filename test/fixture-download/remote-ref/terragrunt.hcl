inputs = {
  name = "World"
}

terraform {
  source = "git::git@github.com:gads-citron/terragrunt.git//test/fixture-download/hello-world?ref=fixture/test"
}
