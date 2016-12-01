.PHONY: build artifact release

all: build

build:
	go build $(FLAGS) -o mtaproxy

artifact:
	GOOS=linux GOARCH=amd64 go build $(FLAGS) -o mtaproxy-linux-amd64

release: artifact
	gzip mtaproxy-linux-amd64
	mv mtaproxy-linux-amd64.gz mtaproxy.gz
	scp mtaproxy.gz do:
