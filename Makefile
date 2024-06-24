PROJECT_NAME=github.com/seedcx/vwap-engine
CONTAINER_NAME=vwap-engine:0.0.1
SRC_MAIN=cmd/
BIN_TARGET=build
BIN_NAME=vwap-engine
TARGET=${BIN_TARGET}

image:
	docker build --build-arg target=${TARGET} -t ${CONTAINER_NAME} .

run:
	docker run -p 3000:3000 ${CONTAINER_NAME}

start: image run

build:
	CGO_ENABLED=0 go build -o ${TARGET}/${BIN_NAME} ${SRC_MAIN}main.go


install: clean build

test:
	go test -v ${PROJECT_NAME}/./...

dep:
	go mod download

vet:
	go vet ${PROJECT_NAME}/./...

lint:
	golangci-lint run --enable-all

clean:
	rm -rf ${TARGET}

clean_build: clean build

cbr: clean build image
