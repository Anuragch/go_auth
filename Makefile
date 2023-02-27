APP?=bin/mainApp
SRC?=cmd/go_auth
clean:
	rm -f ${APP}/*

build: clean
	go build -o ${APP} ${SRC}/main.go

run: build
	 ./${APP}