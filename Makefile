.PHONY: build up down restart logs clean

DOCKER_COMPOSE=docker-compose
SERVICE_NAME?=

build:
	${DOCKER_COMPOSE} build ${SERVICE_NAME}

up:
	${DOCKER_COMPOSE} up -d ${SERVICE_NAME}

down:
	${DOCKER_COMPOSE} down

restart:
ifdef SERVICE_NAME
	${DOCKER_COMPOSE} restart ${SERVICE_NAME}
else
	${DOCKER_COMPOSE} restart
endif

logs:
ifdef SERVICE_NAME
	${DOCKER_COMPOSE} logs -f ${SERVICE_NAME}
else
	${DOCKER_COMPOSE} logs -f
endif

clean:
	${DOCKER_COMPOSE} down -v
	docker system prune -f

logs:
	docker-compose logs -f ${SERVICE_NAME}