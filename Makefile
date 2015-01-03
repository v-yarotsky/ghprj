PREFIX=/usr/local

all: bin/gh install

bin/gh:
	go build -o bin/gh src/github.com/v-yarotsky/main.go
install: bin/gh
	cp -f bin/gh $(PREFIX)/bin/gh
clean:
	rm -f bin/*
