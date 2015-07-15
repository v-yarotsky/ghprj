PREFIX=/usr/local

all: bin/gh alfred install

bin/gh:
	mkdir -p bin
	godep go build -o bin/gh ./main.go

install: bin/gh
	cp -f bin/gh $(PREFIX)/bin/gh

alfred: bin/gh
	rm -rf assets
	mkdir -p assets
	zip -j assets/GithubPrj.alfredworkflow bin/gh alfred_workflow/*
clean:
	rm -f bin/*
	rm -f assets/*
