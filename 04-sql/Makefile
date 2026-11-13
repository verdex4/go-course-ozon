include .env

.PHONY: run all-queries query1 query2 query3 query4 query5 down clean restart psql logs

# Запустить БД
run:
	docker-compose up -d

# Запустить все запросы из queries/
all-queries:
	echo "Запускается запрос №1"
	docker exec -i hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < queries/task_01.sql
	echo "Запускается запрос №2"
	docker exec -i hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < queries/task_02.sql
	echo "Запускается запрос №3"
	docker exec -i hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < queries/task_03.sql
	echo "Запускается запрос №4"
	docker exec -i hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < queries/task_04.sql
	echo "Запускается запрос №5"
	docker exec -i hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < queries/task_05.sql

# Запустить первый запрос (1 задание)
query1:
	docker exec -i hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < queries/task_01.sql

# Запустить второй запрос (2 задание)
query2:
	docker exec -i hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < queries/task_02.sql

# Запустить первый запрос (3 задание)
query3:
	docker exec -i hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < queries/task_03.sql

# Запустить четвертый запрос (4 задание)
query4:
	docker exec -i hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < queries/task_04.sql

# Запустить пятый запрос (5 задание)
query5:
	docker exec -i hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < queries/task_05.sql

# Остановить БД
down:
	docker-compose down

# Очистить БД (удалить volume)
clean:
	docker-compose down -v

# Перезапустить БД (удалить и создать заново)
restart:
	docker-compose down -v
	docker-compose up -d

# Зайти в интерактивный режим psql
psql:
	docker exec -it hw04-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

# Посмотреть логи
logs:
	docker-compose logs -f

