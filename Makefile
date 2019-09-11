build:
	go build .

run: build
	./addon-manager --debug --headless

test:
	go test -v
