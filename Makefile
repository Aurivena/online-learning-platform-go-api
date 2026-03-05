dc = docker-compose

gcs:
	cd resources && goose create $(filter-out $@,$(MAKECMDGOALS)) sql && cd ..
%:
	@:

dc-restart:
	$(dc) down -v && $(dc) up

dc-up:
	$(dc) up
