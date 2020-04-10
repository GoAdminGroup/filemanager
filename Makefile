GOCMD = go

all: fmt

fmt:
	GO111MODULE=off $(GOCMD) fmt ./...