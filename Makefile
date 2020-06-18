
re-build: down build-up

build-up: build up

build:
	docker-compose build --no-cache

up:
	docker-compose up -d

down:
	docker-compose down

push:
	docker-compose push