.PHONY: run generate

run:
	@set -a; \
	if [ -f .env ]; then . ./.env; fi; \
	set +a; \
	go run ./cmd/server

generate:
	buf generate
