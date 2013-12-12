COVER:=$(strip $(shell go tool | grep cover))
ifdef COVER
    COVER:=-cover
endif

all:
	@go get code.google.com/p/goauth2/oauth
	@go get labix.org/v2/mgo
	@go install server

test:
	go test ${COVER} neuroinformatics.harvard.edu/survana
	
