.PHONY: build-web

ifeq ($(OS), Windows_NT)
SRC_FILES := $(shell powershell -Command "(Get-ChildItem -Path database,internal -Recurse).FullName")
else
SRC_FILES := $(shell find database internal)
endif


# This will always build, since it's not really reliable to check if the web files were modified,
# WHILE also checking if they have been built already
build/aio: build-web $(SRC_FILES)
	@echo Building AIO Shipyard
	go build -o build/aio Shipyard/cmd/aio

build-web:
	@echo Compiling the client web interface
	@echo Make sure PNPM is installed
	make -C web

build/remote: $(SRC_FILES)
	@echo Building Shipyard remote
	go build -o build/remote Shipyard/cmd/docker_remote_env