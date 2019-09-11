build:
	go build .

run: build
	./wow-addon-manager --debug

test:
	go test -v
