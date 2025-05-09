CONTAINER_TAG=debian-armv7hf-golang-gtk3-gtk4
CONTAINER_NAME=debian-armv7hf-golang-gtk3-gtk4

build-armv7hf-docker-container:
	docker build -t ${CONTAINER_TAG} -f debian-armv7hf-golang-gtk3-gtk4.Dockerfile .

build-armv7hf-gotk3-example-initial:
	docker run -it --platform linux/arm/v7 -v "$(CURDIR):/workspace" --name ${CONTAINER_NAME} ${CONTAINER_TAG} \
	bash -c "CGO_ENABLED=1 go build -tags 'no_x11' -v -o gotk3-example main.go && exit"

build-armv7hf-gotk3-example:
	docker start ${CONTAINER_NAME} && \
	docker exec -it ${CONTAINER_NAME} bash -c "CGO_ENABLED=1 go build -tags 'no_x11' -v -o gotk3-example main.go && exit" && \
	docker stop ${CONTAINER_NAME}

build-gotk3-example:
	CGO_ENABLED=1 go build -v -o gotk3-example main.go