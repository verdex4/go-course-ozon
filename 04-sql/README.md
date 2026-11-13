## Запуск

Из корня репозитория:
```
cd 04-sql
make run
```

Запустить все запросы:
`make all-queries`

Запустить конкретный запрос (с первого по пятый):
`make query{номер_задания}`

Другие вспомогательные команды:
- Остановить БД: `make down`
- Очистить БД (удалить volume): `make clean`
- Перезапустить БД (удалить и создать заново): `make restart`
- Зайти в интерактивный режим psql: `make psql`
- Посмотреть логи: `make logs`

## Доп. задание

Находится в bonus-task. 
Реализация функции OldestInCity(ctx context.Context, city string) (domain.User, error) находится в bonus-task/storage/repo.go