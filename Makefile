all: lint build test

build:
	go build ./...

install:
	./scripts/make-install.sh

lint:
	golangci-lint run --skip-dirs=repos --disable-all --enable=golint --enable=vet --enable=gofmt ./...
	find . -name '*.go' | xargs gofmt -w -s

fmt:
	./scripts/gofmt.sh

vet:
	go vet ./check ./cmd/... ./download ./handlers ./tools/...
	go vet ./main.go

staticcheck:
	@[ -x "$(shell which staticcheck)" ] || go install honnef.co/go/tools/cmd/staticcheck
	staticcheck ./...

test:
	 go test -cover ./check ./handlers

start:
	 go run main.go

misspell:
	@[ -x "$(shell which misspell)" ] || go install ./vendor/github.com/client9/misspell/cmd/misspell
	find . -name '*.go' -not -path './vendor/*' -not -path './_repos/*' -not -path './download/test_downloads/*' | xargs misspell -error
