Структура проекта

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа, запуск cron планировщика
├── internal/
│   ├── scheduler/
│   │   └── refresh.go          # Логика автоматического обновления данных
│   ├── converter/
│   │   ├── attendance.go        # Конвертер посещаемости (Excel → JSON)
│   │   ├── statement.go         # Конвертер ведомости (Excel → JSON)
│   │   └── xls_to_xlsx.py      # Python скрипт для конвертации XLS → XLSX
│   └── database/
│       └── loader.go            # Загрузка JSON в БД (TODO)
└── go.mod
```

Запуск

Сборка и запуск сервера:

```bash
cd backend
go build ./cmd/server
./server
```

Или напрямую:

```bash
cd backend
go run ./cmd/server
```

Как это работает
1. При старте** сервер сразу запускает обновление данных (конвертирует оба файла)
2. Каждые 90 минут автоматически обновляет данные:
   - `Посещаемость.xlsx` → `public/attendance.json`
   - `ведомость.xls` → `public/summary.json`

Интеграция с существующими скриптами

Старые скрипты (`converter/main.go` и `statement-converter/main.go`) остаются на месте и работают независимо.
Новые функции в `backend/internal/converter/` можно вызывать из любого места в бэкенде.

Зависимости
- `github.com/xuri/excelize/v2` - работа с Excel файлами
- `github.com/robfig/cron/v3` - планировщик задач


Документация:
- API эндпоинты: см. SWAGGER.md
- Swagger UI: http://localhost:8080/swagger/ (после запуска сервера)
- Настройка БД: см. DATABASE.md
