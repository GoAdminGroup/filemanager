GOCMD = go

all: fmt

fmt:
	GO111MODULE=off go fmt ./...
	GO111MODULE=off goimports -l -w .

test:
	gotest -v ./tests