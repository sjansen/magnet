# Bootstrapping

1. `pushd terraform/env/staging`
1. `terragrunt apply -target=module.bootstrap`
1. `popd`
1. Note `ssm-prefix` value.
1. My Security Credentials > CloudFront key pairs > Create New Key Pair
1. Systems Mananager > Parameter Store > fill in values under `ssm-prefix`.
1. `scripts/staging-docker-push`
1. `pushd terraform/env/staging`
1. `terragrunt apply`
