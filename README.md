# go-scheduler

Планировщик HTTP-задач на Go с cron-расписанием, PostgreSQL, Gin и Swagger.

## Запуск

```bash
git clone https://github.com/EthernetUser/go-scheduler.git
cd go-scheduler
docker compose up -d
go run cmd/migrator/main.go
go run cmd/app/main.go
```

## Генерация Swagger

```bash
swag fmt
swag init -g ./cmd/app/main.go -o cmd/docs --parseDependency --parseInternal
```

Swagger: `http://localhost:8080/swagger/index.html`

## API

| Метод  | Путь          | Описание        |
|--------|---------------|-----------------|
| GET    | `/jobs`       | список задач    |
| POST   | `/jobs`       | создать задачу  |
| GET    | `/jobs/{id}`  | получить задачу |
| DELETE | `/jobs/{id}`  | удалить задачу  |

Поля задачи: `name`, `description`, `cron`, `url`.

## Команды Taskfile

```bash
task migrate      # миграции БД
task run          # запуск приложения
task swagger      # генерация Swagger
task dockerup     # поднять Docker
task dockerdown   # остановить Docker
```