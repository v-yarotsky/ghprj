all:
	go build -o bin/gh src/github.com/v-yarotsky/main.go
clean:
	rm -f bin/*
