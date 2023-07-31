PROJECT_NAME := cyclimse/fediverse-blahaj
BINARIES := $(shell find cmd -name '*.go' -type f -exec basename {} \; | sed 's/\.go//g')

.DEFAULT_GOAL := build

JSON_SCHEMA_FILES := $(shell find pkg -name '*.schema.json')
GENERATED_MODELS := $(JSON_SCHEMA_FILES:.schema.json=.schema.json.go)

GENERATED_API := internal/api/v1/api.go

# For all files ending in .schema.json use gojsonschema to generate the model
# Remove the .schema.json suffix from the target name
# Used by the crawler to generate nodeinfo models based on the spec
%.schema.json.go: %.schema.json
	@echo "Generating $@..."
	@gojsonschema -p github.com/$(PROJECT_NAME)/$(shell dirname $<) --resolve-extension schema.json -o $@ $<

# Generates the server boilerplate from the OpenAPI spec
$(GENERATED_API): api/v1/openapi.yaml
	@echo "Generating $@..."
	@oapi-codegen -package v1 \
		-generate "types,server,strict-server,spec" \
		$< > $@

.PHONY: build
build: $(GENERATED_MODELS) $(GENERATED_API)
	for binary in $(BINARIES); do \
		echo "Building $$binary..."; \
		go build -o bin/$$binary ./cmd/$$binary/$$binary.go; \
	done

.PHONY: test
test: $(GENERATED_MODELS)
	@echo "Testing..."
	@go test ./...

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f $(GENERATED_MODELS) $(BIN_NAME)