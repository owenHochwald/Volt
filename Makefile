APP_EXECUTABLE=volt

build:
	go build -o ${APP_EXECUTABLE} ./cmd/volt/main.go

build-mac:
	GOARCH=amd64 GOOS=darwin go build -o ${APP_EXECUTABLE} ./cmd/volt/main.go

run: build
	./${APP_EXECUTABLE}

clean:
	go clean
	rm -f ${APP_EXECUTABLE}

.PHONY: build build-mac run clean