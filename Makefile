GO ?= go
GOFMT_FILES ?= $$(find . -name '*.go' | grep -v vendor)
TF_LOG ?= error

KEA_DHCP4_ADDRESS ?= http://localhost:8000
KEA_DHCP4_HTTP_USERNAME ?=
KEA_DHCP4_HTTP_PASSWORD ?=

default: build

build:
	$(GO) build -v

test:
	$(GO) get -t ./...
	$(GO) test -parallel $$(nproc) -race -timeout 60m -v ./...

testacc:
	KEA_DHCP4_ADDRESS=$(KEA_DHCP4_ADDRESS) \
	KEA_DHCP4_HTTP_USERNAME=$(KEA_DHCP4_HTTP_USERNAME) \
	KEA_DHCP4_HTTP_PASSWORD=$(KEA_DHCP4_HTTP_PASSWORD) \
	TF_LOG=$(TF_LOG) TF_ACC=1 \
	$(GO) test -parallel 4 -v -race $(TESTARGS) -timeout 60m ./internal/...

generate:
	cd tools && $(GO) generate

vet:
	@echo "$(GO) vet ."
	@$(GO) vet $$($(GO) list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@echo "==> Checking that code complies with gofmt requirements..." ; \
	files=$$(find . -name '*.go') ; \
	gofmt_files=$$(gofmt -l $$files); \
	if [ -n "$$gofmt_files" ]; then \
		echo 'gofmt needs running on the following files:'; \
		echo "$$gofmt_files"; \
		echo "You can use the command: \`make fmt\` to reformat code."; \
		exit 1; \
	fi

static-analysis:
	@if command -v golangci-lint > /dev/null; then \
		echo "==> Running golangci-lint"; \
		golangci-lint run --timeout 5m; \
	else \
		echo "Missing \"golangci-lint\" command, not linting .go" >&2; exit 1; \
	fi
	@if command -v terraform > /dev/null; then \
		echo "==> Running terraform fmt"; \
		terraform fmt -recursive -check -diff; \
	elif command -v tofu > /dev/null; then \
		echo "==> Running tofu fmt"; \
		tofu fmt -recursive -check -diff; \
	else \
		echo "Missing \"terraform\" command, not checking .tf format" >&2; exit 1; \
	fi

update-gomod:
	$(GO) get -t -v -u ./...
	$(GO) mod tidy -go=1.25
	$(GO) get toolchain@none
	@echo "Dependencies updated"

.PHONY: build test testacc vet fmt fmtcheck generate static-analysis update-gomod
