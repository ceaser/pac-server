.PHONY: build run stop bash clean distclean push pull

DOCKER_APPNAME=pac-server
VERSION=$(shell git describe --tag 2>/dev/null | (grep -Eo "[0-9\.\+\:\~]+-?[0-9]+" | head -n1 || echo 0.0.0))
BUILD_OBJS=$(DOCKER_APPNAME)_$(VERSION)_amd64.deb
ARGS=-p 80:80

build: $(BUILD_OBJS)
	docker build --build-arg PKG=$(BUILD_OBJS) -t "$(DOCKER_APPNAME):$(VERSION)" .
	rm $(BUILD_OBJS)

run:
	docker run -d $(ARGS) --name="$(DOCKER_APPNAME)" "$(DOCKER_APPNAME):$(VERSION)"

debug: build
	docker run -it --entrypoint /bin/bash $(ARGS) --rm $(DOCKER_APPNAME):$(VERSION)

bash:
	docker exec -it "${DOCKER_APPNAME}" bash

clean:
	-docker stop ${DOCKER_APPNAME}
	-docker rm -v ${DOCKER_APPNAME}

distclean: clean
	-rm $(BUILD_OBJS)
	-rm -rf mnt

push:
	docker tag "$(DOCKER_APPNAME):$(VERSION)" "$(DOCKER_APPNAME):latest"
	docker push "$(DOCKER_APPNAME):$(VERSION)"
	docker push "$(DOCKER_APPNAME):latest"

pull:
	docker pull "$(DOCKER_APPNAME):$(VERSION)"

$(BUILD_OBJS):
	cp ../*.deb .
