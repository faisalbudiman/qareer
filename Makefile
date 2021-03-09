VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
NAME="qareer/api"
PORT=8000
DOCKER_URL="180266367152.dkr.ecr.ap-southeast-1.amazonaws.com"
ifeq ($(BRANCH), HEAD)
	BRANCH=$(TRAVIS_BRANCH)
endif

.PHONY: build
build:  
	@export DOCKER_CONTENT_TRUST=1 && docker build --build-arg=VERSION=${VERSION} --build-arg=BRANCH=${BRANCH} --build-arg=BUILD=${BUILD} --build-arg=NAME=${NAME} --build-arg=PORT=${PORT} -f api.dockerfile -t ${NAME} .

.PHONY: build-no-cache
build-no-cache: 
	@export DOCKER_CONTENT_TRUST=1 && docker build --no-cache --build-arg=VERSION=${VERSION} --build-arg=BRANCH=${BRANCH} --build-arg=BUILD=${BUILD} --build-arg=NAME=${NAME} --build-arg=PORT=${PORT} -f api.dockerfile -t ${NAME} .


.PHONY: run
run:  
	@docker run --rm -p ${PORT}:${PORT} -e PORT=':${PORT}' ${NAME}:latest 

.PHONY: tag-stag
tag-stag:
	@docker tag ${NAME}:latest ${DOCKER_URL}/${NAME}:stag

.PHONY: tag-prod
tag-prod:
	@docker tag ${NAME}:latest ${DOCKER_URL}/${NAME}:prod

.PHONY: tag-version
tag-version:
	@docker tag ${NAME}:latest ${DOCKER_URL}/${NAME}:${VERSION}

.PHONY: push
push:
	@docker push ${DOCKER_URL}/${NAME}
