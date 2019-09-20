build:
	go build -o wow-addon-manager .

run: build
	go run .

test:
	go test -v
