.PHONY: build install alfred clean

PREFIX=/usr/local

all: build alfred install

build:
	mkdir -p bin
	rm -f bin/gh
	godep go build -o bin/gh ./main.go

install: build
	cp -f bin/gh $(PREFIX)/bin/gh

alfred: build
	rm -rf assets
	mkdir -p assets
	zip -j assets/GithubPrj.alfredworkflow bin/gh alfred_workflow/*
clean:
	rm -f bin/*
	rm -f assets/*
