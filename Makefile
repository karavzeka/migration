version := 1.0

.PHONY: bin/*

param_1 := $(word 2,$(MAKECMDGOALS))
param_2 := $(word 3,$(MAKECMDGOALS))
param_3 := $(word 4,$(MAKECMDGOALS))
# Used in 'create' target
ifneq ($(param_2),)
	migration_name := $(param_2)
else
	migration_name := "noname"
endif

network_name := migration-network
image_name := pgmycli

create: CMD = create
create: chk_param_1 recompile image
	@docker run --rm --name migration-create-cmd --mount type=bind,source=$(PWD),target=/var/project $(image_name):$(version) bin/create --dir /var/project/databases/$(param_1)/migrations --name $(migration_name)

version: CMD = version
version: chk_param_1 network recompile image
	@docker run --user `id -u`:`id -g` --rm --name migration-version-cmd --network $(network_name) --mount type=bind,source=$(PWD),target=/var/project $(image_name):$(version) bin/version $(param_1) $(param_2) $(param_3)

init: CMD = init
init: chk_param_1 network recompile image
	@docker run --user `id -u`:`id -g` --rm --name migration-init-cmd --network $(network_name) --mount type=bind,source=$(PWD),target=/var/project $(image_name):$(version) bin/init $(param_1) $(param_2)

migrate: CMD = migrate
migrate: chk_param_1 chk_param_2 network recompile image
	@docker run --user `id -u`:`id -g` --rm --name migration-migrate-cmd --network $(network_name) --mount type=bind,source=$(PWD),target=/var/project $(image_name):$(version) bin/migrate $(param_1) $(param_2) $(param_3)

snapshot: CMD = snapshot
snapshot: chk_param_1 network recompile image
	@docker run --user `id -u`:`id -g` --rm --name migration-snapshot-cmd --network $(network_name) --mount type=bind,source=$(PWD),target=/var/project $(image_name):$(version) bin/snapshot $(param_1)

network:
	@docker network inspect $(network_name) > /dev/null 2>&1 || (echo 'Network "$(network_name)" has created' && docker network create $(network_name))

image:
	@docker image inspect $(image_name):$(version) > /dev/null 2>&1 || docker build -t $(image_name):$(version) .

recompile:
	@test bin/$(CMD) -nt src/cmd/$(CMD) -a bin/$(CMD) -nt "$$(find src/internal -type f -exec stat --format '%Y :%n' "{}" \; | sort -nr |  head -n 1 | cut -d: -f2-)" || \
	docker run --rm --user `id -u`:`id -g` --mount type=bind,source=$(PWD),target=/go/project --env GOCACHE=/go/cache golang sh -c "cd project/src; go build -o ../bin/$(CMD) cmd/$(CMD)/main.go"

chk_param_1:
ifeq ($(param_1),)
	@echo "You have an error in command. Look at the example:"
	@$(MAKE) --no-print-directory $(CMD)_help
	@echo ""
	@false
endif

chk_param_2:
ifeq ($(param_2),)
	@echo "You have an error in command. Look at the example:"
	@$(MAKE) --no-print-directory $(CMD)_help
	@echo ""
	@false
endif

create_help:
	@echo "make create <db_alias> [<migration_tag>]"

version_help:
	@echo "make version <db_alias> [force] [<version>]"
	@echo "    force   - set database version and clean it"
	@echo "    version - define version, using with force"

init_help:
	@echo "make init <db_alias>"

migrate_help:
	@echo "make migrate <db_alias> <cmd> [<version>]"
	@echo "    <cmd> - up|goto|down"

snapshot_help:
	@echo "make snapshot <db_alias>"

%:
	@true