.phony: up down ps logs

up:
	docker-compose up -d

down:
	docker-compose down

ps:
	docker ps

logs:
	docker-compose logs -f

re: down up

generate:
	go generate ./...
