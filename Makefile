.PHONY: up
up:
	@echo docker-compose up in background
	@docker-compose up -d
	@docker network connect divalocal diva-challenge

.PHONY: build
build:
	@echo Build new image with tag
	@docker build -t planningdiva/diva-challenge .
.PHONY: down
down:
	@echo docker-compose down
	@docker-compose down

.PHONY: rebuild
rebuild:
	@echo Docker-compose down
	@docker-compose down
	@echo Remove old image
	@docker rmi planningdiva/diva-challenge
	@echo Build new image with tag
	@docker build -t planningdiva/diva-challenge .
	@echo docker-compose up in background
	@docker-compose up -d
	@docker network connect divalocal diva-challenge
