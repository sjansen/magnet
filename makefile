.PHONY:  default  check-env  check-working-tree
.PHONY:  login  staging  snapshot  start  test


default: start


check-env:
ifndef AWSCLI
	$(error AWSCLI is undefined)
endif



check-working-tree:
	@git diff-index --quiet HEAD -- \
	|| (echo "Working tree is dirty. Commit all changes."; false)


dynamodb:
	@scripts/docker-up-dynamodb


login: check-env
	$(AWSCLI) ecr get-login-password \
	| docker login \
	    --username AWS \
	    --password-stdin \
	    `scripts/get-ecr-registry`


snapshot:
	git archive -o snapshot.tgz HEAD


staging: check-working-tree login
	scripts/staging-docker-push
	scripts/staging-deploy


start:
	foreman start -e /dev/null


test:
	@scripts/docker-up-test
	@echo ' ____'
	@echo '|  _ \ __ _ ___ ___ '
	@echo '| |_) / _` / __/ __|'
	@echo '|  __/ (_| \__ \__ \'
	@echo '|_|   \__,_|___/___/'
