build:
	go build -o wow-addon-manager .

run:
	go run . --fast

run-slow:
	go run .

test:
	go test -v

ci:
	nix-shell shell.nix --run 'xvfb-run make test run-slow'
