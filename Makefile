all:
	@export GOPATH=$(CURDIR)
	@echo GOPATH: ${GOPATH}
	@go get code.google.com/p/goauth2/oauth
	@go get labix.org/v2/mgo
	@go install server

test:
	@export GOPATH=$(CURDIR)
	@echo GOPATH: ${GOPATH}
	go test --coverprofile=cover.out neuroinformatics.harvard.edu/survana
