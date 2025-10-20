#!/bin/bash

# Wait for Debezium Connect to be ready
echo "Waiting for Debezium Connect to start..."
until curl -s -o /dev/null -w "%{http_code}" http://localhost:8083/connectors | grep -q "200"; do
    echo "Debezium Connect is not ready yet. Waiting..."
    sleep 5
done

echo "Debezium Connect is ready!"

# Check if connector already exists
CONNECTOR_EXISTS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8083/connectors/postgres-connector)

if [ "$CONNECTOR_EXISTS" = "200" ]; then
    echo "Connector 'postgres-connector' already exists. Deleting it first..."
    curl -X DELETE http://localhost:8083/connectors/postgres-connector
    sleep 2
fi

# Register the PostgreSQL connector
echo "Registering PostgreSQL connector..."
curl -X POST http://localhost:8083/connectors \
    -H "Content-Type: application/json" \
    -d @postgres-connector.json

echo ""
echo "Connector registered successfully!"
echo ""

# Check connector status
echo "Connector status:"
curl -s http://localhost:8083/connectors/postgres-connector/status | jq '.'

echo ""
echo "Available topics in Kafka:"
docker exec kafka kafka-topics --list --bootstrap-server localhost:9092 | grep -E "(products|customers|orders|order_items)"

echo ""
echo "Setup complete!"
echo ""
echo "You can now:"
echo "1. Access Kafka UI at http://localhost:8080"
echo "2. View Debezium connectors at http://localhost:8083/connectors"
echo "3. Connect to PostgreSQL at localhost:5432 (user: postgres, password: postgres, database: testdb)"
echo ""
echo "To see messages in a topic, run:"
echo "docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic products --from-beginning"