img_tag = "hayeah/qtumportal:testjanus"
container_name = "qtumportal-testjanus"

.PHONY: renewJanus
renewJanus:
	GOOS=linux GOARCH=amd64 go build -o ./docker/janus/janus-linux-amd64 github.com/dcb9/janus/cli/janus
	cd docker \
		&& docker-compose stop janus  \
		&& docker-compose rm -f janus  \
		&& docker-compose build janus \
		&& docker-compose up -d janus \
		&& cd ../

.PHONY: up
up:
	cd docker \
		&& docker-compose up -d 

.PHONY: stop
stop:
	cd docker \
		&& docker-compose stop

.PHONY: init-container
init-container:
	docker exec -it qtumportal-testjanus sh init-container.sh

.PHONY: logs
logs:
	cd docker \
		&& docker-compose logs -f
