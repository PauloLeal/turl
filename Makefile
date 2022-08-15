TARGET_FILE:=${shell head -n1 go.mod | sed -r 's/.*\/(.*)/\1/g' }
BUILD_DIR=.build
COVER_PROFILE_FILE="${BUILD_DIR}/go-cover.tmp"
RESOURCES_FILE=resources/resources.go
RESOURCES_FILE_EXISTS=$(shell [ -e $(RESOURCES_FILE) ] && echo 1 || echo 0 )

.PHONY: target clean mk-build-dir update-deps build-deps b build build-all clean-test test cover-html badge

target: build

clean:
	@rm -rf $(TARGET_FILE) $(BUILD_DIR)

############## build tasks

mk-build-dir:
	@mkdir -p ${BUILD_DIR}

update-deps:
	@go get -u -d -v ./...
	
build-deps:
	@go get -d -v ./...

b: 
	@go build -o $(TARGET_FILE) cmd/main.go

build: clean build-deps test
	@$(MAKE) b

build-all: clean mk-build-dir build-deps test
	GOOS=linux $(MAKE) b && zip -9 $(TARGET_FILE)-linux64.zip $(TARGET_FILE) && rm $(TARGET_FILE)
	GOOS=windows $(MAKE) b && mv $(TARGET_FILE) $(TARGET_FILE).exe && zip -9 $(TARGET_FILE)-win64.zip $(TARGET_FILE).exe && rm $(TARGET_FILE).exe
	GOOS=darwin $(MAKE) b && zip -9 $(TARGET_FILE)-osx64.zip $(TARGET_FILE) && rm $(TARGET_FILE)
	@mv *.zip ${BUILD_DIR}
	@find $(BUILD_DIR) -type f

############## test tasks

clean-test:
	@go fmt ./...
	@go clean -testcache

test: clean-test
	go test -p 1 ./...

cover-html: mk-build-dir clean-test
	go test -p 1 -coverprofile=${COVER_PROFILE_FILE} ./...
	go tool cover -html=${COVER_PROFILE_FILE}
	$(MAKE) badge

badge:
	@go install github.com/jpoles1/gopherbadger@latest
	gopherbadger -md="README.md" -png=false 1>&2 2> /dev/null
	@if [ -f coverage.out ]; then rm coverage.out ; fi; 
