PB = $(wildcard *.proto)
MODULE = $(PB:.proto=)

all: $(MODULE)

# Using `make block` to generate block.pb.go
%: %.proto
		protoc -I=. -I=$(GOPATH)/src --gogofaster_out=plugins=grpc:$(GOPATH)/src $<

clean:
		rm -f ./*.pb.go

.PHONY: help
help:
	@echo  ''
	@echo  '!Make sure the project is in your GOPATH: $(GOPATH) '
	@echo  ''
	@echo  'Declare go_package option in *.proto file, and import proto file with GOPATH.'
	@echo  'For example, `block.proto` needs import `transaction.proto`, write proto file like follow:'
	@echo '# syntax = "proto3";'
	@echo '# package protos;'
	@echo '#'
	@echo '# option go_package = "github.com/hyperchain/go-hpc-common/types/protos";'
	@echo '#'
	@echo '# import "github.com/hyperchain/go-hpc-common/types/protos/transaction.proto";'
	@echo '#'
	@echo '# ...'
	@echo ''
	@echo '> Using `make block` to generate block.pb.go'
	@echo '> Using `make all` to generate all'
	@echo '> Using `make clean` to remove all generated files'
	@echo ''
