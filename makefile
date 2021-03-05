.PHONY:  default  check-env  check-working-tree  login
.PHONY:  push-staging  snapshot  start  test


default: start


check-env:
ifndef AWSCLI
	$(error AWSCLI is undefined)
endif



dynamodb:
	@scripts/docker-up-dynamodb


login: check-env
	$(AWSCLI) ecr get-login-password \
	| docker login \
	    --username AWS \
	    --password-stdin \
	    `scripts/get-ecr-registry`


push-staging: login
	@scripts/push-staging


snapshot:
	git archive -o snapshot.tgz HEAD


start:
	foreman start -e /dev/null


test:
	@scripts/docker-up-test
	@echo ' ____'
	@echo '|  _ \ __ _ ___ ___ '
	@echo '| |_) / _` / __/ __|'
	@echo '|  __/ (_| \__ \__ \'
	@echo '|_|   \__,_|___/___/'
