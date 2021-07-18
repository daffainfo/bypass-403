fmt:
	@gofmt -w -s main.go && goimports -w main.go && go vet main.go