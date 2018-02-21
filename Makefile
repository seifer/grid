SERVICE_PORT=8080
DOCKER_SERVICE_IMAGE=grid

run:
	@docker run -it --rm -p $(SERVICE_PORT):$(SERVICE_PORT) $(DOCKER_SERVICE_IMAGE):latest

test:
	@curl -s -XPOST --data '@./data/test_request.json' http://127.0.0.1:$(SERVICE_PORT)/fetch

build:
	@docker build --rm --no-cache \
		-t $(DOCKER_SERVICE_IMAGE):latest \
		-t $(DOCKER_SERVICE_IMAGE):1.0 .