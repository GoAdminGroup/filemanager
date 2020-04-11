GOCMD = go

all: fmt

fmt:
	GO111MODULE=off $(GOCMD) fmt ./...

test:
	gotest -v ./tests