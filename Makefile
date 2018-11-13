
all:
	go build

test:
	go test

cover:
	go test -coverprofile cover.out

show:
	go tool cover -html=cover.out

