BUILD_DIR = ./bin
BINARY_NAME = crypto-server
BUILD_FLAGS := -mod=vendor -v -ldflags "-w -s"

clean:
	@rm -rf $(BUILD_DIR)

build: clean
	@mkdir -p $(BUILD_DIR) > /dev/null
	go build ${BUILD_FLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ./main.go
	@echo "binary created at ${BUILD_DIR}/${BINARY_NAME}"