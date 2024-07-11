# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/gitfix
BINARY_NAME := gitfix
INSTALL_DIR := /usr/local/bin
ZSHRC_FILE := ~/.zshrc
LOCAL_MAN := ./manpage/gitfix.8
MAN_DIR := /opt/homebrew/share/man/man8
# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
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

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -o ${MAIN_PACKAGE_PATH}/${BINARY_NAME}

# ==================================================================================== #
# INSTALLATION
# ==================================================================================== #
install: build move_bin move_man

# move binary to wherever
.PHONY: move_bin
move_bin:
	@echo "Installing $(BINARY_NAME)..."
	chmod +x $(MAIN_PACKAGE_PATH)/$(BINARY_NAME)
	sudo cp $(MAIN_PACKAGE_PATH)/$(BINARY_NAME) $(INSTALL_DIR)
	chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "Updating .zshrc to include $(INSTALL_DIR)..."

.PHONY: move_man
move_man:
	@echo "Installing gitfix manpage"
	sudo cp ./manpage/gitfix.8 $(MAN_DIR)