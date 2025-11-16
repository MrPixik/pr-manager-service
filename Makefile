up_local: # запуск бд и миграций через докер. Само приложение запустится через go run
	docker-compose -f docker-compose.local.yaml up -d
	go run cmd/main.go --env=local

down_local: # остановка контейнеров. чтобы завершить работу сервиса необходимо еще отправить ctr+c в консоль
	docker compose -f docker-compose.local.yaml stop

run_tests:
	go test -v ./internal/http/server/handlers/user
	go test -v ./internal/http/server/handlers/team
	go test -v ./internal/http/server/handlers/pull_request
	go test -v ./internal/service/user
	go test -v ./internal/service/team
	go test -v ./internal/service/pull_request

up_prod: # запуск всего сервера (подгружается образ с моего dockerHub)
	docker-compose -f docker-compose.prod.yaml up -d

down_prod: # остановка всех контейнеров
	docker compose -f docker-compose.prod.yaml stop pr-manager-service
	docker compose -f docker-compose.prod.yaml stop migrations
	docker compose -f docker-compose.prod.yaml stop postgres

deploy_local: #для личного удобства
	docker build -t pr-manager-service:latest -f ./deploy/docker/Dockerfile .

deploy_push_remote: #для личного удобства
	docker login
	docker tag pr-manager-service:latest pixik/pr-manager-service:latest
	docker push pixik/pr-manager-service:latest