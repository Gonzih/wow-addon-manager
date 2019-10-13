build:
	go build -o wow-addon-manager .

run:
	go run . --fast

test:
	go test -v
