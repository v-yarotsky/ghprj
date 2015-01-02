PREFIX=/usr/local

all:
	go build -o bin/gh src/github.com/v-yarotsky/main.go
install: all
	cp -f bin/gh $(PREFIX)/bin/gh
clean:
	rm -f bin/*
