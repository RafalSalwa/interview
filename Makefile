.PHONY:

build:
	docker compose up --build -d
up:
	docker compose up -d --force-recreate && docker compose logs -f gateway auth_service user_service consumer_service crypton_service
compose-down:
	docker compose down --remove-orphans

logs_follow:
	docker compose logs -f gateway auth_service user_service consumer_service crypton_service

.PHONY: tester
tester:
	docker compose up -f docker-compose.tester.yml -d

test_unit:
	APP_ENV=staging go test -v -cover -coverprofile=cover.out ./pkg/... ./cmd/... -tags=unit
	go tool cover -html=cover.out -o coverage.html

test_integration:
	APP_ENV=staging go test -cover ./cmd/... -tags=integration

lint:
	golangci-lint run --enable gosec

check_sec:
	gosec ./...

static_analysis:
	goimports
	errcheck ./...
	gofmt -s
	go vet ./...
	staticcheck ./...

.PHONY: proto
proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
		protoc --proto_path=proto --go_out=proto/grpc --go_opt=paths=source_relative   --go-grpc_out=proto/grpc --go-grpc_opt=paths=source_relative   proto/*.proto; \

clean:
	go clean -i google.golang.org/grpc/...
