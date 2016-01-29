.PHONY: cmd

PREFIX=/usr/local

all: cmd alfred install

cmd:
	mkdir -p bin
	go build -o bin/gh cmd/main.go

install: cmd
	cp -f bin/gh $(PREFIX)/bin/gh

alfred: cmd
	rm -rf assets
	mkdir -p assets
	zip -j assets/GithubPrj.alfredworkflow bin/gh alfred_workflow/*
clean:
	rm -f bin/*
	rm -f assets/*
