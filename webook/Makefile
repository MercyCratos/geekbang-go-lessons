.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=amd64 go build -tags=k8s -o webook .
	@docker rmi -f mercycrazy/webook:v0.0.1
	@docker build -t mercycrazy/webook:v0.0.1 .