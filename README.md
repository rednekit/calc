# Распределённый вычислитель арифметических выражений

## Обзор

Этот проект реализует распределённую систему для вычисления арифметических выражений. Он состоит из двух основных компонентов:

1. **Оркестратор**: Принимает арифметические выражения, преобразует их в набор задач и обеспечивает порядок их выполнения.
2. **Агент**: Получает задачи от оркестратора, выполняет их и возвращает результаты.

## Возможности

- Поддержка операций сложения, вычитания, умножения и деления.
- Каждая операция может выполняться независимо и параллельно.
- Масштабируемость за счёт добавления новых агентов для обработки большего числа задач.
- Предоставляет API для отправки выражений, проверки статуса вычислений и получения результатов.

## Начало работы

### Предварительные требования

- Go 1.18 или новее

### Переменные окружения

Установите следующие переменные окружения для настройки времени выполнения операций в миллисекундах:

- `TIME_ADDITION_MS`
- `TIME_SUBTRACTION_MS`
- `TIME_MULTIPLICATION_MS`
- `TIME_DIVISION_MS`

Пример:
```sh
export TIME_ADDITION_MS=1000
export TIME_SUBTRACTION_MS=1000
export TIME_MULTIPLICATION_MS=1000
export TIME_DIVISION_MS=1000
```

## Сборка и запуск

### Запуск оркестратора:
cd calc
go run server.go


### Запуск агента:
cd calc/agent
go run main.go

## API Эндпоинты
Оркестратор

### Добавить выражение
POST /api/v1/calculate
Content-Type: application/json

{
    "expression": "2 + 2 * 2"
}

### Получить все выражения
GET /api/v1/expressions

### Получить выражение по ID
GET /api/v1/expressions/{id}

### Получить задачу
GET /internal/task

### Отправить результат задачи
POST /internal/task
Content-Type: application/json

{
    "id": 1,
    "result": 4.0
}

## Примеры

### Добавить выражение
curl --location --request POST 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "expression": "12 * 17 + 18 * 22"
}'

Windows syntax:
curl -H "Content-Type: application/json" -X POST http://localhost:8080/api/v1/calculate -d "{\"expression\": \"12 * 17 + 18 * 22\"}"

### Получить все выражения
curl --location 'localhost:8080/api/v1/expressions'

### Получить выражение по ID
curl --location 'localhost:8080/api/v1/expressions/1'



Этот `README.md` файл предоставляет описание проекта, инструкции по настройке и запуску, а также примеры использования API на русском языке.
