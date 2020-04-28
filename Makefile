GOCMD=go
GOFMT1=gofmt
GOFMT2=goreturns
GOFMT3=goimports
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GODOC=$(GOCMD) doc
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

all: fmt build doc
ci: build test bench
doc:
	$(GODOC) -all .
fmt:
	$(GOFMT1) -s -w .
	$(GOFMT2) -l -w .
	$(GOFMT3) -w .
	$(GOMOD) tidy
build:
	$(GOBUILD) -v
test:
	$(GOTEST) -v -race -short -cover -covermode=atomic
bench:
	$(GOTEST) -parallel=4 -run="none" -benchtime="2s" -benchmem -bench=.
