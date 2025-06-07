run: build
	./app

build:
	go build -o app

run-prod:
	@env $$(kamal secrets print) $(MAKE) run

test:
	go test ./...

context:
	uvx files-to-prompt --extension go .
