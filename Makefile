.PHONY: run
run:
	go run cmd/main.go

.PHONY: swag
swag:
	swag init -g cmd/main.go -pd

.PHONY: web
web:
	cd ./web && npm run build:prod


.PHONY: image
image:
	docker build . -t  goadmin:latest


.PHONY: deploy
deploy:
	cd ./deploy &&  docker compose up -d