appName=pac-server
version=$(shell git describe --tags 2> /dev/null || git rev-parse HEAD)
BUILD_OBJS=$(appName)

build:
	go build -ldflags "-X github.com/ceaser/pac/internal/version.AppName=${appName} -X github.com/ceaser/pac/internal/version.Version=${version} -X github.com/ceaser/pac/internal/version.BuildTimeUTC=`date -u '+%Y-%m-%d_%I:%M:%S%p'`"

install: build
	go install

clean:
	-rm -rf $(BUILD_OBJS) 2>/dev/null 1>&2

#distclean: clean
#	-$(MAKE) -C docker distclean
