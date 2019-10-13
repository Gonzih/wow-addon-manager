build:
	go build -o wow-addon-manager .

run:
	go run . --fast

test:
	go test -v

ci:
	nix-shell shell.nix --run 'make test run'
