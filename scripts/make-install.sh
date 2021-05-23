#!/bin/sh

go get github.com/golangci/golangci-lint/cmd/golangci-lint
go install ./vendor/github.com/client9/misspell/cmd/misspell
