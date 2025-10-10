

all:
	go build
	./patron ./test/*.pat

test: all
	cd test && go test

install: all
	cp patron ~/opt/bin
