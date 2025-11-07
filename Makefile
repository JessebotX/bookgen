GOEXE = go
GOFMTEXE = gofmt
GOFMTFLAGS = -s -w -l
GOFMTINPUT = $(SRC)
GOLINTEXE = staticcheck
GOLINTFLAGS =
GOLINTINPUT = ./...

SRC = *.go cmd/pub/*.go renderer/html/*.go
BIN = pub
TESTDATAPATH = testdata/book1

all: build

build:
	$(GOEXE) mod download
	$(GOEXE) build ./cmd/$(BIN)

fmt:
	$(GOFMTEXE) $(GOFMTFLAGS) $(GOFMTINPUT)

lint:
	$(GOLINTEXE) $(GOLINTFLAGS) $(GOLINTINPUT)

clean: clean-bin clean-test

clean-bin:
	rm -f $(BIN) $(BIN).exe

clean-test:
	rm -rf $(TESTDATAPATH)/_output
