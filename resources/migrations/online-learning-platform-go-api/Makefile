dc = docker-compose

gc_create:
	cd resources/migrations && goose create $(filter-out $@,$(MAKECMDGOALS)) sql && cd ../..
%:
	@:

gc_up:
	cd resources/migrations && goose up && cd ../..

dc-restart:
	$(dc) down -v && $(dc) up

dc-up:
	$(dc) up
