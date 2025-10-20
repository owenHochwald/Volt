APP_EXECUTABLE=volt

build:
	GOARCH=amd64 GOOS=darwin go build -o ${APP_EXECUTABLE} ./cmd/volt/main.go

run: build
	./${APP_EXECUTABLE}

clean:
	go clean
	rm ${APP_EXECUTABLE}
