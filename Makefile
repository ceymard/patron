

all:
	go build
	./patron ./test/*.pat

test: all
	cd test && go test

install:
	cp patron ~/opt/bin
