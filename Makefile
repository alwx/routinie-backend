VERSION = 8
IMAGE_NAME = alwxdev/routinie-backend:$(VERSION)

development-deps-up:
	docker-compose -f deployments/docker-compose.development.yml up --build

development-deps-down:
	docker-compose -f deployments/docker-compose.development.yml down

build-static:
	cp -r static/public/* static/dist

build-backend: build-static
	docker build --platform linux/amd64 --rm -t $(IMAGE_NAME) -f Dockerfile .
	docker push $(IMAGE_NAME)

build-all: build-backend

run: build-static
	go generate cmd/backend/main.go
	go run cmd/backend/main.go