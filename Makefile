
all:
	go build -o markdown-format github.com/cofyc/markdown-format/cmd/markdown-format

test:
	go test -timeout 5m github.com/cofyc/markdown-format/...
