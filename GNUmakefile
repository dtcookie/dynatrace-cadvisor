BUILD_DIR := build

default: build

all: build

clean:
ifeq ($(OS),Windows_NT)
	@if exist "$(BUILD_DIR)" rmdir /S /Q $(BUILD_DIR)
else
	@rm -Rf $(BUILD_DIR)
endif

build: export GOOS=linux
build: export GOARCH=amd64
build:
	go build -o $(BUILD_DIR)/dynatrace-cadvisor

release: clean build
ifdef TAG
	docker rmi -f dtcookie/dynatrace-cadvisor:$(TAG)
	docker build -t dtcookie/dynatrace-cadvisor:$(TAG) .
	docker push dtcookie/dynatrace-cadvisor:$(TAG)
endif

test: 
	go test -v