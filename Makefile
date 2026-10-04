App=url-shortener
GEN_PATH="internal/pkg/server/gen"
GEN_TOOL="./tools/openapi-generator-cli.jar"
OPEN_API_GENERATOR_IGNORE="./api/.openapi-generator-ignore"
URL_SHORTENER_OPENAPI_SPEC="api/url-shortener/url-shortener.openapi.yaml"

build:
	go build -o bin/$(APP) cmd/url_shortener/main.go

run:
	go run cmd/url_shortener/main.go

test:
	go test ./...

lint:
	golangci-lint run

docker:
	docker build -t $(APP) .

generate-server:
	@if [ -d $(App) ]; then \
		echo "Removing existing generated server code..."; \
		rm -r $(GEN_PATH); \
	fi
	@if [ ! -d "internal/pkg/server" ]; then \
		echo "Target directory missing. Initializing workspace layout..."; \
		mkdir -p internal/pkg/server; \
	fi
	@cp $(OPEN_API_GENERATOR_IGNORE) internal/pkg/server/; \
	java -jar $(GEN_TOOL) generate \
		-i $(URL_SHORTENER_OPENAPI_SPEC) \
		-g go-server \
		-o internal/pkg/server/ \
		-c server-config.yml
	
	@echo "Cleaning up generator metadata artifacts..."
	@rm -rf internal/pkg/server/.openapi-generator
	rm  internal/pkg/server/.openapi-generator-ignore
