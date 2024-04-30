# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## build: build containers and services
build:
	docker compose up --build -d
up:
	docker compose up -d --force-recreate && docker compose logs -f gateway auth_service user_service consumer_service
compose-down:
	docker compose down --remove-orphans

logs_follow:
	docker compose logs -f gateway auth_service user_service consumer_service

.PHONY: tester
tester:
	docker compose up -f docker-compose.tester.yml -d

test_unit:
	APP_ENV=staging gotestsum -- -race -covermode=atomic -cover -coverprofile=coverage.out ./pkg/... -tags=unit -json > coverage.json
	go tool cover -html=coverage.out -o coverage.html

test_unit_cli:
	APP_ENV=staging gotestsum -- -race -covermode=atomic -cover -coverprofile=coverage.out ./pkg/... -tags=unit
	go tool cover -func coverage.out

test_integration:
	APP_ENV=staging go test -cover ./cmd/... -tags=integration

vet:
	go vet ./...

vet_with_json:
	go vet -json ./... > go_vet.json || true

lint:
	golangci-lint run

lint_with_checkstyle:
	golangci-lint run --out-format checkstyle > golangci-lint.xml || true

check_sec:
	gosec ./...

static_analysis:
	goimports
	errcheck ./...
	gofmt -s
	go vet ./...
	staticcheck ./...

.PHONY: sonar_static_analysis
sonar_static_analysis: test_unit
	go vet -json ./... > go_vet.json
	golangci-lint run --out-format checkstyle > golangci-lint.xml
	sonar-scanner -Dsonar.host.url=${SONAR_HOST} -Dsonar.token=${SONAR_TOKEN}

check_static:
	#go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...
check_callvis:
	#go install github.com/ofabry/go-callvis@latest
	go-callvis github.com/RafalSalwa/auth-api/cmd/gateway
check_goreporter:
	#go get -u github.com/360EntSecGroup-Skylar/goreporter
	
check_revive:
	#go install github.com/mgechev/revive@latest
	revive -config revive.toml -formatter unix ./...

check_review_dog:
	#go install github.com/reviewdog/reviewdog/cmd/reviewdog@latest

.PHONY: proto
proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
		protoc --proto_path=proto --go_out=proto/grpc --go_opt=paths=source_relative   --go-grpc_out=proto/grpc --go-grpc_opt=paths=source_relative   proto/*.proto; \

clean:
	go clean -i google.golang.org/grpc/...
