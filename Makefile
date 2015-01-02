all:
	go build -o bin/gh src/github.com/v-yarotsky/*
clean:
	rm -f bin/*
