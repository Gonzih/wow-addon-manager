build:
	go build -o wow-addon-manager .

run:
	go run . --fast

run-slow:
	go run .

test:
	go test -v

ci:
	xvfb-run make test run-slow

deps:
	sudo add-apt-repository ppa:longsleep/golang-backports -y
	sudo apt update
	sudo apt install -y golang-go xvfb chromium-browser git
