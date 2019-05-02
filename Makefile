.PHONY: clean build test acceptance all
GO_SOURCES = $(shell find . -type f -name '*.go')

all: test build acceptance

build: artifactory/io/projectriff/java/io.projectriff.java

test:
	go test -v ./...

acceptance:
	pack create-builder -b acceptance/testdata/builder.toml projectriff/builder
	docker pull packs/run:0.1.0
	GO111MODULE=on go test -v -tags=acceptance ./acceptance

artifactory/io/projectriff/java/io.projectriff.java: buildpack.toml $(GO_SOURCES)
	rm -fR $@ 							&& \
	./ci/package.sh						&& \
	mkdir $@/latest 					&& \
	tar -C $@/latest -xzf $@/*/*.tgz

clean:
	rm -fR artifactory/
	rm -fR dependency-cache/