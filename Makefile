all:
	@go get code.google.com/p/goauth2/oauth
	@go get labix.org/v2/mgo
	@go install server

test:
	@export GOPATH=${HOME}/gopath/src/${REPOSITORY}
	@echo GOPATH: ${GOPATH}
	go test --coverprofile=cover.out neuroinformatics.harvard.edu/survana
