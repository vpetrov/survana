BIN_DIR:=bin
WWW_DIR:=www
SSL_DIR:=${BIN_DIR}/ssl
BUILD_DIR:=build
DIST_DIR:=dist

TARGET:=${BIN_DIR}/server
SURVANA:=${BIN_DIR}/survana
SERVER_CONF:=${BIN_DIR}/survana.json

MAKEFILE_PATH:= $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR:= $(dir $(MAKEFILE_PATH))

COVER:=$(strip $(shell go tool | grep cover))
ifdef COVER
    COVER:=-cover
endif

COVERFILE := survana.coverage

ifdef DEBUG
	GO_FLAGS:=-gcflags "-N -l"
endif

#OSX
OSX_PROJECT_NAME:=Survana
OSX_TARGET:=bin/${OSX_PROJECT_NAME}.app
OSX_BUILD_ARCHIVE:=build/${OSX_PROJECT_NAME}.xcarchive
OSX_BUILD_RESOURCES:=${OSX_BUILD_ARCHIVE}/Products/Applications/Survana.app/Contents/Resources/


OSX_DIST_APP:=${OSX_PROJECT_NAME}.app
OSX_DIST_APP_PATH:=${DIST_DIR}/${OSX_DIST_APP}
OSX_DIST_TARGET:=survana-osx.zip
OSX_DIST_TARGET_PATH:=${DIST_DIR}/${OSX_DIST_TARGET}


OSX_ROOT:=src/platform/osx
OSX_PROJECT_ROOT:=${OSX_ROOT}/${OSX_PROJECT_NAME}
OSX_PROJECT:=${OSX_PROJECT_ROOT}/${OSX_PROJECT_NAME}.xcodeproj

XCODE:=xcodebuild

GIT_VERSION:=$(strip $(shell git describe --tags))

all: ${TARGET} ${SURVANA}

${SURVANA}: src/survana/download.go
	GOPATH=${CURRENT_DIR} go install ${GO_FLAGS} survana

${TARGET}:
	GOPATH=${CURRENT_DIR} go get code.google.com/p/goauth2/oauth
	GOPATH=${CURRENT_DIR} go get labix.org/v2/mgo
	GOPATH=${CURRENT_DIR} go get github.com/vpetrov/perfect
	GOPATH=${CURRENT_DIR} go get github.com/vpetrov/perfect/auth
	GOPATH=${CURRENT_DIR} go install ${GO_FLAGS} server

test:
	GOPATH=${CURRENT_DIR} go test ${COVER} neuroinformatics.harvard.edu/survana

clean:
	@rm -f ${TARGET}
	GOPATH=${CURRENT_DIR} go clean -i neuroinformatics.harvard.edu/survana
	GOPATH=${CURRENT_DIR} go clean -i neuroinformatics.harvard.edu/survana/dashboard
	GOPATH=${CURRENT_DIR} go clean -i neuroinformatics.harvard.edu/survana/study
	GOPATH=${CURRENT_DIR} go clean -i neuroinformatics.harvard.edu/survana/store

distclean:
	@rm -rf ${TARGET} ${DIST_DIR}
	GOPATH=${CURRENT_DIR} go clean -i neuroinformatics.harvard.edu/survana
	GOPATH=${CURRENT_DIR} go clean -i neuroinformatics.harvard.edu/survana/dashboard
	GOPATH=${CURRENT_DIR} go clean -i neuroinformatics.harvard.edu/survana/study
	GOPATH=${CURRENT_DIR} go clean -i neuroinformatics.harvard.edu/survana/store
	GOPATH=${CURRENT_DIR} go clean -i github.com/vpetrov/perfect
	GOPATH=${CURRENT_DIR} go clean -i github.com/vpetrov/perfect/auth
	GOPATH=${CURRENT_DIR} go clean -i labix.org/v2/mgo
	GOPATH=${CURRENT_DIR} go clean -i labix.org/v2/mgo/bson
	GOPATH=${CURRENT_DIR} go clean -i code.google.com/p/goauth2/oauth
	
osx: ${TARGET} ${OSX_TARGET}

osx-dist: osx
	@mkdir -p ${DIST_DIR}
	@echo "[Cleaning]"
	@rm -rf ${OSX_DIST_APP_PATH} ${OSX_DIST_TARGET_PATH}
	@echo "[Copying app]"
	@cp -r ${OSX_TARGET} ${OSX_DIST_APP_PATH}
	@echo "[Creating archive]"
	@cd ${DIST_DIR} && \
	 zip -9 -q -r ${OSX_DIST_TARGET} ${OSX_DIST_APP} && \
	 cd ${CURRENT_DIR}
	@echo "Done. See ${OSX_DIST_TARGET_PATH}"

${OSX_TARGET}: ${OSX_BUILD_ARCHIVE} ${OSX_PROJECT}
	@rm -rf ${OSX_TARGET}
	${XCODE} -project ${OSX_PROJECT} -exportArchive -exportFormat APP -archivePath ${OSX_BUILD_ARCHIVE} -exportPath ${OSX_TARGET}
	
${OSX_BUILD_ARCHIVE}: ${OSX_PROJECT}
	${XCODE} -project ${OSX_PROJECT} -scheme ${OSX_PROJECT_NAME} archive -archivePath ${OSX_BUILD_ARCHIVE}
	@mkdir -p ${OSX_BUILD_RESOURCES}/services/server
	@cp -r ${BIN_DIR}/server ${BIN_DIR}/ssl ${OSX_BUILD_RESOURCES}/services/server/

osx-clean:
	rm -rf ${OSX_TARGET} ${OSX_BUILD_ARCHIVE}

cover: ${TARGET}
	GOPATH=${CURRENT_DIR} go test -coverprofile=${COVERFILE} -covermode=count neuroinformatics.harvard.edu/survana
	GOPATH=${CURRENT_DIR} go tool cover -html=${COVERFILE}

bench: ${TARGET}
	GOPATH=${CURRENT_DIR} go test -bench . neuroinformatics.harvard.edu/survana

format: ${TARGET}
	GOPATH=${CURRENT_DIR} go fmt neuroinformatics.harvard.edu/survana
