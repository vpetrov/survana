BIN_DIR:=bin
WWW_DIR:=www
SSL_DIR:=${BIN_DIR}/ssl

TARGET:=${BIN_DIR}/server
SERVER_CONF:=${BIN_DIR}/survana.json

MAKEFILE_PATH:= $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR:= $(dir $(MAKEFILE_PATH))

COVER:=$(strip $(shell go tool | grep cover))
ifdef COVER
    COVER:=-cover
endif

#OSX
OSX_PROJECT_NAME:=Survana
OSX_TARGET:=bin/${OSX_PROJECT_NAME}.app
OSX_BUILD_ARCHIVE:=build/${OSX_PROJECT_NAME}.xcarchive
OSX_ARCHIVE_SERVER_DIR:=${OSX_BUILD_ARCHIVE}/Products/Applications/${OSX_PROJECT_NAME}.app/Contents/Resources/server

OSX_ROOT:=src/platform/osx
OSX_PROJECT_ROOT:=${OSX_ROOT}/${OSX_PROJECT_NAME}
OSX_PROJECT:=${OSX_PROJECT_ROOT}/${OSX_PROJECT_NAME}.xcodeproj

XCODE:=xcodebuild

all: ${TARGET}

${TARGET}:
	GOPATH=${CURRENT_DIR} go get code.google.com/p/goauth2/oauth
	GOPATH=${CURRENT_DIR} go get labix.org/v2/mgo
	GOPATH=${CURRENT_DIR} go get github.com/coopernurse/gorp
	GOPATH=${CURRENT_DIR} go get github.com/mattn/go-sqlite3 
	GOPATH=${CURRENT_DIR} go install server

test:
	go test ${COVER} neuroinformatics.harvard.edu/survana

clean:
	@rm -f ${TARGET}
	
osx: ${TARGET} ${OSX_TARGET}

${OSX_TARGET}: ${OSX_BUILD_ARCHIVE} ${OSX_PROJECT}
	cp -r ${BIN_DIR}/server ${SERVER_CONF} ${WWW_DIR} ${SSL_DIR} ${OSX_ARCHIVE_SERVER_DIR}
	${XCODE} -project ${OSX_PROJECT} -exportArchive -exportFormat APP -archivePath ${OSX_BUILD_ARCHIVE} -exportPath ${OSX_TARGET}
	
${OSX_BUILD_ARCHIVE}: ${OSX_PROJECT}
	${XCODE} -project ${OSX_PROJECT} -scheme ${OSX_PROJECT_NAME} archive -archivePath ${OSX_BUILD_ARCHIVE}


osx-clean:
	rm -rf ${OSX_TARGET} ${OSX_BUILD_ARCHIVE}
