# Machine-specific settings are read from the environment.
# Put them in a .env file (see .env.example) or export them in your shell.
-include .env
export

SWAG_VERSION := v1.8.9
SWAGGER_URL ?= http://localhost:9000/docs/index.html

.DEFAULT_GOAL := dev

.PHONY: dev run-backend-dev run-frontend-dev open-frontend open-swagger gen-docs install-swagger check-frontend-env

## Run backend, frontend and open the browser
dev:
	$(MAKE) -j 3 run-backend-dev run-frontend-dev open-frontend

## Run the API in dev mode with hot reload (air)
run-backend-dev:
	air

## Run the frontend dev server (needs FRONTEND_DIR)
run-frontend-dev: check-frontend-env
	sleep 2 && cd "$(FRONTEND_DIR)" && yarn dev

## Open the frontend in the browser (needs FRONTEND_URL)
open-frontend: check-frontend-env
	sleep 3 && open "$(FRONTEND_URL)"

## Open the Swagger UI in the browser (override with SWAGGER_URL)
open-swagger:
	open "$(SWAGGER_URL)"

## Regenerate Swagger docs into ./docs
gen-docs:
	swag init -g cmd/main.go

## Install the swag CLI version the docs are generated with
install-swagger:
	go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)

check-frontend-env:
	@test -n "$(FRONTEND_DIR)" || (echo "FRONTEND_DIR is not set (see .env.example)" && exit 1)
	@test -n "$(FRONTEND_URL)" || (echo "FRONTEND_URL is not set (see .env.example)" && exit 1)
