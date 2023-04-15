SWAG = $(shell which swag)

build: swagger buildApp

swagger: ## Runs 'swagger' in the local build environment
	@${SWAG} init -g main.go -o swagger; \
    [ -e ./swagger/doc.json ] && rm -v ./swagger/doc.json; \
    [ -s ./swagger/swagger.json ] && mv -v ./swagger/swagger.json ./swagger/doc.json

buildApp:
	go mod tidy
	go build -o NORSI-TRANS

