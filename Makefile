# Variables
BINARY_NAME=app
BUILD_DIR=build
GO_BUILD_CMD=go build -o ${BUILD_DIR}/${BINARY_NAME} cmd/app/main.go


# Targets
.PHONY: build run clean

build:clean
	${GO_BUILD_CMD}

run: build
	./${BUILD_DIR}/${BINARY_NAME} -threads=5

clean:
	rm ${BINARY_NAME}
