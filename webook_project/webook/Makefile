.PHONY: docker
docker:
	@rm webook || true
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f darling/webook:v0.0.1 .
	@docker build -t darling/webook:v0.0.1 .


#.PHONY: docker
#docker:
#	@rm webook || true
#	@$env:CGO_ENABLED=0
#	@$env:GOOS="linux"
#	@$env:GOARCH="arm"
#	@go build -o webook .
#	@docker rmi -f darling/webook:v0.0.1 .
#	@docker build -t darling/webook:v0.0.1 .