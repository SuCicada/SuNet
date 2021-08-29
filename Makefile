
export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build:
	bash ./build.sh

fmt:
	@sh -c "go fmt ./..."

wc:
	@echo result count: bin/
	@ls -lR ./bin | grep ^-.* | wc -l

clean:
	@rm -rf ./bin/*
	echo clean over
