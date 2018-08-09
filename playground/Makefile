img_tag = "hayeah/qtumportal:testjanus"
container_name = "qtumportal-testjanus"

.PHONY: image
image:
	GOOS=linux GOARCH=amd64 go build -o ./docker/janus-linux-amd64 github.com/dcb9/janus/cli/janus
	docker build ./docker -t $(img_tag)

.PHONY: run
run:
	cp -n ./docker/env.list-example ./docker/env.list || true
	docker run --rm -d \
	  --name $(container_name) \
	  --env-file ./docker/env.list \
	  -v `pwd`/docker:/dapp \
	  -p 3889:3889 \
	  -p 23889:23889 \
	  $(img_tag)
	make init-container

.PHONY: stop
stop:
	docker stop $(container_name)

.PHONY: exec
exec:
	docker exec -it $(container_name) sh

.PHONY: init-container
init-container:
	docker exec -it qtumportal-testjanus sh init-container.sh

.PHONY: logs
logs:
	docker logs --tail 1000  -f $(container_name)

.PHONY: clean-qtum-data
clean-qtum-data:
	rm -rf ./docker/.qtum
