test:
	go test -v ./... -cover

clean:
	find . -name flymake_* -delete

sloccount:
	 find . -path ./Godeps -prune -o -name "*.go" -print0 | xargs -0 wc -l
