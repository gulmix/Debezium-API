# Debezium PostgreSQL to Kafka Setup

Полная конфигурация для запуска Debezium с PostgreSQL, Kafka и Kafka-UI.

## Компоненты

- **PostgreSQL** - База данных с тестовыми таблицами и данными
- **Apache Kafka** - Брокер сообщений
- **Zookeeper** - Координатор для Kafka
- **Schema Registry** - Хранилище схем
- **Debezium Connect** - CDC коннектор для PostgreSQL
- **Kafka-UI** - Веб-интерфейс для управления Kafka

## Быстрый старт

1. Запустите все сервисы:
```bash
docker-compose up -d
```

2. Дождитесь запуска всех контейнеров (около 30 секунд):
```bash
docker-compose ps
```

3. Зарегистрируйте Debezium коннектор:
```bash
./setup-connector.sh
```

## Доступ к сервисам

- **Kafka-UI**: http://localhost:8080
- **Debezium Connect API**: http://localhost:8083
- **PostgreSQL**: localhost:5432
  - User: postgres
  - Password: postgres
  - Database: testdb
- **Schema Registry**: http://localhost:8081

## Тестовые данные

База данных содержит схему `inventory` с таблицами:
- `products` - Товары (10 записей)
- `customers` - Покупатели (10 записей)
- `orders` - Заказы (10 записей)
- `order_items` - Позиции заказов

## Работа с данными

### Просмотр сообщений в топике
```bash
# Просмотр всех сообщений о продуктах
docker exec -it kafka kafka-console-consumer \
    --bootstrap-server localhost:9092 \
    --topic products \
    --from-beginning

# Просмотр новых сообщений о заказах
docker exec -it kafka kafka-console-consumer \
    --bootstrap-server localhost:9092 \
    --topic orders
```

### Проверка статуса коннектора
```bash
curl http://localhost:8083/connectors/postgres-connector/status | jq
```

### Список топиков
```bash
docker exec kafka kafka-topics --list --bootstrap-server localhost:9092
```

### Тестирование CDC

Подключитесь к PostgreSQL и измените данные:
```bash
docker exec -it postgres psql -U postgres -d testdb

-- Обновить цену продукта
UPDATE inventory.products SET price = 999.99 WHERE id = 1;

-- Добавить нового покупателя
INSERT INTO inventory.customers (first_name, last_name, email, phone)
VALUES ('Test', 'User', 'test@email.com', '555-9999');

-- Создать новый заказ
INSERT INTO inventory.orders (customer_id, status, total_amount)
VALUES (1, 'NEW', 100.00);
```

Изменения автоматически появятся в соответствующих Kafka топиках.

## Остановка и очистка

```bash
# Остановить все контейнеры
docker-compose down

# Полная очистка с удалением volumes
docker-compose down -v
```

## Конфигурация коннектора

Файл `postgres-connector.json` содержит настройки Debezium коннектора:
- Отслеживание всех таблиц в схеме `inventory`
- Использование логической репликации PostgreSQL (pgoutput)
- Начальный снимок всех существующих данных
- Преобразование имен топиков для удобства

## Troubleshooting

### Проверка логов
```bash
# Логи Debezium
docker logs debezium

# Логи PostgreSQL
docker logs postgres

# Логи Kafka
docker logs kafka
```

### Перезапуск коннектора
```bash
# Удалить коннектор
curl -X DELETE http://localhost:8083/connectors/postgres-connector

# Заново создать
./setup-connector.sh
```