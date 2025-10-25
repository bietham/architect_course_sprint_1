# temperature-api-go

Простой Go-сервис, который отдаёт случайную температуру.

## Сборка и запуск (Docker Compose)
```bash
docker compose up --build
```

Доступные сервисы:
- temperature-api: http://localhost:8081
- smart_home-db: Postgres (db=smart_home, user=smart, pass=smart_pass)

## Эндпоинты
- `GET /healthz` — проверка жизни
- `GET /temperature?location=&sensorId=` — возвращает JSON со случайной температурой.
  - Если `location` отсутствует, берётся значение из `sensorId` (1=Living Room, 2=Bedroom, 3=Kitchen).
  - Если `sensorId` отсутствует — вычисляется на основе `location`.

Примеры:
```bash
curl "http://localhost:8081/temperature?location=Living%20Room"
curl "http://localhost:8081/temperature?sensorId=2"
curl "http://localhost:8081/temperature"
```
