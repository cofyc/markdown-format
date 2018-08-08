
all:
	go build -o markdown-toc github.com/cofyc/markdown-toc/cmd/markdown-toc

test:
	go test -timeout 5m github.com/cofyc/markdown-toc/...
