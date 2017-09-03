appName=pac
version=$(shell git describe --tags 2> /dev/null || git rev-parse HEAD)
debVersion = $(shell git describe --tag 2>/dev/null | (grep -Eo "[0-9\.\+\:\~]+-?[0-9]+" | head -n1 || echo 0.0.0))
BUILD_OBJS=$(appName) $(appName)_$(debVersion)_amd64.deb

build:
	go build -ldflags "-X github.com/ceaser/pac/internal/version.AppName=${appName} -X github.com/ceaser/pac/internal/version.Version=${version} -X github.com/ceaser/pac/internal/version.BuildTimeUTC=`date -u '+%Y-%m-%d_%I:%M:%S%p'`"

deb: clean
	GOOS=linux GOARCH=amd64 $(MAKE) build
	fpm -t deb \
		-s dir \
		-n $(appName) \
		--version $(debVersion) \
		--license "Copyright (c) 2016-$(shell date +%Y) Ceaser Larry. All rights reserved" \
		--vendor "Ceaser Larry" \
		--url "http://ca.utio.us" \
		--category misc \
		--deb-priority optional \
		--maintainer "Ceaser Larry <c@utio.us>" \
		--architecture amd64 \
		--description "PAC configuration" \
		--deb-changelog debian/changelog \
		--after-install debian/$(appName).postinst \
		--deb-init debian/$(appName).init  \
		--deb-default debian/$(appName).default \
		--package $(appName)_$(debVersion)_amd64.deb \
		$(appName)=/usr/sbin/$(appName) \

install: build
	mv ${appName} $$GOBIN/

clean:
	rm -rf $(BUILD_OBJS)

