ENGINE=docker

.PHONY: build-docker
build-docker:
	docker run -it --rm -v $(shell pwd):/app -p 8000:8000 amzani/docker-gatsbyjs build

.PHONY: dev-docker
dev-docker:
	@@echo "Deleting cache"
	rm -r ./.cache/* || true
	@@echo "Run in docker"
	docker run -it --rm -v $(shell pwd):/app -p 8000:8000 amzani/docker-gatsbyjs develop

.PHONY: build
build:
	rm -fr ./public/*
	ACTIVE_ENV=dev gatsby build

.PHONY: dev
dev:
	@@echo "Deleting cache"
	rm -r ./.cache/* || true
	gatsby develop

.PHONY: serve-dev
serve-dev: build
	@@echo "Running in docker with dev server"
	docker-compose down --rmi local
	docker-compose up