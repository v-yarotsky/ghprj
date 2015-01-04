PREFIX=/usr/local

all: bin/gh alfred install

bin/gh:
	mkdir -p bin
	go build -o bin/gh src/github.com/v-yarotsky/main.go

install: bin/gh
	cp -f bin/gh $(PREFIX)/bin/gh

alfred: bin/gh
	rm -rf assets
	mkdir -p assets
	zip -j assets/GithubPrj.alfredworkflow bin/gh alfred_workflow/*
clean:
	rm -f bin/*
	rm -f assets/*
