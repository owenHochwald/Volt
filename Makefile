APP_EXECUTABLE=volt
INSTALL_PATH=/usr/local/bin

build:
	go build -o ${APP_EXECUTABLE} ./cmd/volt/main.go

build-mac:
	GOARCH=amd64 GOOS=darwin go build -o ${APP_EXECUTABLE} ./cmd/volt/main.go

install: build
	sudo cp ${APP_EXECUTABLE} ${INSTALL_PATH}/

uninstall:
	sudo rm -f ${INSTALL_PATH}/${APP_EXECUTABLE}

run: build
	./${APP_EXECUTABLE}

clean:
	go clean
	rm -f ${APP_EXECUTABLE}

.PHONY: build build-mac install uninstall run clean