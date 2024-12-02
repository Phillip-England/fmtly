install:
	go build -o gtml main.go; rm -r /home/jacex/.local/bin/gtml; mv gtml /home/jacex/.local/bin;

test-all:
	go test -run TestAll

test:
	go test -run TestComponents