.PHONY: build run stop bash clean distclean push pull

DOCKER_REPO=ceaser
DOCKER_APPNAME=pac-server
VERSION=$(shell git describe --tag 2>/dev/null | (grep -Eo "[0-9\.\+\:\~]+-?[0-9]+" | head -n1 || echo 0.0.0))
BUILD_OBJS=$(DOCKER_APPNAME)_$(VERSION)_amd64.deb
ARGS=-p 80:80

build: $(BUILD_OBJS)
	docker build --build-arg PKG=$(BUILD_OBJS) -t "$(DOCKER_REPO)/$(DOCKER_APPNAME):$(VERSION)" .
	rm $(BUILD_OBJS)

run:
	docker run -d $(ARGS) --name="$(DOCKER_APPNAME)" "$(DOCKER_REPO)/$(DOCKER_APPNAME):$(VERSION)"

debug: build
	docker run -it --entrypoint /bin/bash $(ARGS) --rm $(DOCKER_REPO)/$(DOCKER_APPNAME):$(VERSION)

bash:
	docker exec -it "${DOCKER_APPNAME}" bash

clean:
	-docker stop ${DOCKER_APPNAME} 2>/dev/null 1>&2
	-docker rm -v ${DOCKER_APPNAME} 2>/dev/null 1>&2

distclean: clean
	-rm $(BUILD_OBJS) 2> /dev/null
	-docker rmi $(DOCKER_REPO)/$(DOCKER_APPNAME):$(VERSION) 2>/dev/null 1>&2
	-docker images | grep  "<none>" | awk '{ print $$3 }' | xargs docker rmi 2>/dev/null 1>&2

push:
	docker tag "$(DOCKER_REPO)/$(DOCKER_APPNAME):$(VERSION)" "$(DOCKER_REPO)/$(DOCKER_APPNAME):latest"
	docker push "$(DOCKER_REPO)/$(DOCKER_APPNAME):latest"

pull:
	docker pull "$(DOCKER_REPO)/$(DOCKER_APPNAME):latest"

$(BUILD_OBJS):
	cp ../$(BUILD_OBJS) .
