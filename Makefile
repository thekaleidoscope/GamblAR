.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	# @dep ensure
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@docker-compose up --force-recreate -d
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@docker-compose down
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./GamblAR

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/GamblAR-* GamblAR
	@docker rm -f -v `docker ps -a --no-trunc | grep "GamblAR" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "GamblAR" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"