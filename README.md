# go-cloud-backend

## Запуск

Требования:
- Go `1.26.2`

1. Установите зависимости:

```bash
go mod download
```

2. Создайте файл `.env` в корне проекта:

```bash
cp .env.example .env
```

3. Запустите API:

```bash
go run ./cmd/api
```

4. Проверьте, что сервис работает:

```bash
curl http://localhost:8080/test
```
