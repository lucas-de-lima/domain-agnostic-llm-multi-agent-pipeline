# Domain-agnostic LLM multi-agent pipeline

BINARY_NAME := pipeline
IMAGE_NAME := video-forge
OUTPUT_DIR := output

.PHONY: build run docker-build docker-run clean

# Build the pipeline binary (default: cmd/pipeline)
build:
	go build -o $(BINARY_NAME) ./cmd/pipeline/

# Build the legacy distiller binary
build-distiller:
	go build -o distiller ./cmd/distiller/

# Run locally (requires LLM endpoint via ENV). Usage: make run URL=https://...
run: build
	@if [ -z "$(URL)" ]; then echo "Usage: make run URL=<source>"; exit 1; fi
	mkdir -p $(OUTPUT_DIR)/temp
	./$(BINARY_NAME) "$(URL)"

# Run with specialization. Usage: make run URL=... SPEC=recipes
run-spec: build
	@if [ -z "$(URL)" ]; then echo "Usage: make run-spec URL=<source> [SPEC=recipes]"; exit 1; fi
	mkdir -p $(OUTPUT_DIR)/temp
	./$(BINARY_NAME) -specialization=$(or $(SPEC),default) "$(URL)"

# Build Docker image
docker-build:
	docker build -t $(IMAGE_NAME) .

# Run via Docker. Usage: make docker-run URL=https://...
docker-run: docker-build
	@if [ -z "$(URL)" ]; then echo "Usage: make docker-run URL=<source>"; exit 1; fi
	mkdir -p $(OUTPUT_DIR)/temp
	docker run --rm \
		-v "$(PWD)/$(OUTPUT_DIR):/app/output" \
		--add-host=host.docker.internal:host-gateway \
		$(IMAGE_NAME) "$(URL)"

clean:
	rm -f $(BINARY_NAME) distiller
