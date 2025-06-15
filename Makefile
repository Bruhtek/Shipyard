build-aio: build/aio build-web
	@echo Building AIO Shipyard
	go build -o build/aio Shipyard/cmd/aio

build-web: web/build
	@echo Compiling the client web interface
	@echo Make sure PNPM is installed
	make -C web

build-remote:
	@echo Building Shipyard remote
	go build -o build/remote Shipyard/cmd/docker_remote_env