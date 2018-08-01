APP=product-api

.PHONY: build
build:
	docker build . -t $(APP)


.PHONY: run
run:
	docker run -it --rm --name product-api -p 5000:5000  product-api