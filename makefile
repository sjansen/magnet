.PHONY:  default  check-env  check-working-tree  login  push-staging  run  test


default: run


check-env:
ifndef AWSCLI
	$(error AWSCLI is undefined)
endif


check-working-tree:
	@git diff-index --quiet HEAD -- \
	|| (echo "Working tree is dirty. Commit all changes."; false)


login: check-env
	$(AWSCLI) ecr get-login-password \
	| docker login \
	    --username AWS \
	    --password-stdin \
	    `scripts/get-ecr-registry`


push-staging: check-working-tree login
	@scripts/push-staging


run:
	@scripts/docker-up-dynamodb


test:
	@scripts/docker-up-test
	@echo ' ____'
	@echo '|  _ \ __ _ ___ ___ '
	@echo '| |_) / _` / __/ __|'
	@echo '|  __/ (_| \__ \__ \'
	@echo '|_|   \__,_|___/___/'
