.PHONY: cmd test

all: test cmd alfred

test:
	go test

cmd:
	mkdir -p bin
	go build -o bin/gh cmd/gh/main.go
	go build -o bin/ghlogin cmd/ghlogin/main.go

alfred: cmd
	rm -rf assets
	mkdir -p assets
	zip -j assets/GithubPrj_Alfred3.alfredworkflow bin/gh bin/ghlogin alfred_workflow/*
clean:
	rm -f bin/*
	rm -f assets/*
