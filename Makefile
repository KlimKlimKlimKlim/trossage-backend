.PHONY: mocks
mocks:
	@mockery

.PHONY: swagger
swagger:
	@swag init --v3.1 -g cmd/main.go -o docs

.PHONY: test
test:
	@go test ./...

.PHONY: lint-all
lint-all:
	@./scripts/lint-all.sh $(PWD)

.PHONY: lint-file
lint-file:
	@./scripts/lint-file.sh $(PWD) $(FILE)


.PHONY: docker-up
docker-up:
	@docker compose down --rmi local --remove-orphans
	@docker compose up -d --build --force-recreate --remove-orphans

.PHONY: docker-down
docker-down:
	@docker compose down -v --rmi local --remove-orphans
	@docker image prune -f

.PHONY: docker-deep-clean
docker-deep-clean:
	@docker compose down -v --rmi local --remove-orphans 2>/dev/null || true
	@docker image prune -af
	@docker builder prune -af
	@docker system prune -af --volumes

