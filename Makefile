build:
	go build -o wow-addon-manager .

run: build
	./wow-addon-manager --debug

test:
	go test -v
