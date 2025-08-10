.PHONY: build-web

ifeq ($(OS), Windows_NT)
SRC_FILES := $(shell powershell -Command "(Get-ChildItem -Path src,static -Recurse).FullName")
else
SRC_FILES := $(shell find src static)
endif


build-web: build/index.html

build/index.html: $(SRC_FILES) $(STATIC_FILES)
	pnpm i
	pnpm run build