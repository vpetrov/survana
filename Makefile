BIN_DIR:=bin
WWW_DIR:=www
SSL_DIR:=${BIN_DIR}/ssl
SERVER_CONF:=${BIN_DIR}/survana.json

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

all:
	@go get code.google.com/p/goauth2/oauth
	@go get labix.org/v2/mgo
	@go install server

test:
	go test ${COVER} neuroinformatics.harvard.edu/survana
	
osx: ${OSX_TARGET}

${OSX_TARGET}: ${OSX_BUILD_ARCHIVE} ${OSX_PROJECT}
	cp -r ${BIN_DIR}/server ${SERVER_CONF} ${WWW_DIR} ${SSL_DIR} ${OSX_ARCHIVE_SERVER_DIR}
	${XCODE} -project ${OSX_PROJECT} -exportArchive -exportFormat APP -archivePath ${OSX_BUILD_ARCHIVE} -exportPath ${OSX_TARGET}
	
${OSX_BUILD_ARCHIVE}: ${OSX_PROJECT}
	${XCODE} -project ${OSX_PROJECT} -scheme ${OSX_PROJECT_NAME} archive -archivePath ${OSX_BUILD_ARCHIVE}


osx-clean:
	rm -rf ${OSX_TARGET} ${OSX_BUILD_ARCHIVE}
