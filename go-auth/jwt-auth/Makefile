COMPOSE = nerdctl compose
PROJECT_NAME = go-crud-auth

db-nerdctl-up:
	$(COMPOSE) -p $(PROJECT_NAME) -f compose.yml up -d

db-nerdctl-down:
	$(COMPOSE) -p $(PROJECT_NAME) -f compose.yml down

db-nerdctl-logs:
	$(COMPOSE) -p $(PROJECT_NAME) -f compose.yml logs -f

db-nerdctl-restart:
	$(MAKE) db-nerdctl-down
	$(MAKE) db-nerdctl-up

db-nerdctl-status:
	$(COMPOSE) -p $(PROJECT_NAME) -f compose.yml ps
