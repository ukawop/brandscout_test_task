# BrandScout Test Task - Сервис цитат

In-Memory REST API для управления цитатами.

## Функционал

- Добавление новых цитат
- Просмотр всех цитат
- Получение случайной цитаты
- Фильтрация по автору
- Удаление цитат по ID

## API Endpoints

- `POST /quotes` - Добавить новую цитату
- `GET /quotes` - Получить все цитаты (фильтр по автору `?author=имя`)
- `GET /quotes/random` - Получить случайную цитату
- `DELETE /quotes/{id}` - Удалить цитату по ID

## Запуск проекта

### Требования

- Go версии 1.18 или новее

### Установка

1. Клонировать репозиторий:
   ```bash
   git clone https://github.com/ukawop/brandscout_test_task.git
   cd brandscout_test_task
   ```
2. Запустить сервис:
   ```bash
   make run
   ```
### Проверить работу
1. Добавить новую цитату:
   ```bash
   curl -X POST http://localhost:8080/quotes \
   -H "Content-Type: application/json" \
   -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}' 
   ```
2. Получить все цитаты:
   ```bash
   curl http://localhost:8080/quotes
   ```
3. Получить все цитаты по фильтру:
   ```bash
   curl http://localhost:8080/quotes?author=Confucius
   ```
4. Получить случайную цитату:
   ```bash
   curl http://localhost:8080/quotes/random
   ```
5. Удалить цитату по ID:
   ```bash
   curl -X DELETE http://localhost:8080/quotes/1
   ```
### Запуск тестов
   ```bash
   make test
   ```
