.PHONY: build-dev build-prod run-dev run-prod gen-swagger
build-dev:
	docker build --target development -t recipe-catalog-dev .

build-prod:
	docker build --target production -t recipe-catalog .

run-dev:
	docker run -p 8080:8080 recipe-catalog-dev

run-prod:
	docker run -p 8080:8080 recipe-catalog

gen-swagger:
	swag init --parseDependency --parseInternal -g ./cmd/main.go -o ./docs