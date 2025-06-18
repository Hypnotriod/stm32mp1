IMAGE_TAG=debian-armv7hf-golang-gtk3-gtk4
CONTAINER_NAME=debian-armv7hf-golang-gtk3-gtk4

build-armv7hf-docker-image:
	docker build --platform linux/arm/v7 -t ${IMAGE_TAG} -f debian-armv7hf-golang-gtk3-gtk4.Dockerfile .

build-armv7hf-gotk3-example:
	docker run --rm -it --platform linux/arm/v7 \
	-v "$(CURDIR):/workspace" -v "$(CURDIR)/bin/.gomodcache:/go/pkg/mod" -v "$(CURDIR)/bin/.gocache:/root/.cache/go-build" \
	--name ${CONTAINER_NAME} ${IMAGE_TAG} \
	bash -c "CGO_ENABLED=1 go build -tags 'no_x11' -v -o ./bin/gotk3-example-armv7hf main.go && exit"

build-gotk3-example:
	CGO_ENABLED=1 go build -v -o ./bin/gotk3-example main.go