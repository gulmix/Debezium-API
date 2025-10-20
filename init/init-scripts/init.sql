-- Create test schema
CREATE SCHEMA IF NOT EXISTS inventory;

-- Create products table
CREATE TABLE inventory.products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2),
    quantity INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create customers table
CREATE TABLE inventory.customers (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create orders table
CREATE TABLE inventory.orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES inventory.customers(id),
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'PENDING',
    total_amount DECIMAL(10, 2)
);

-- Create order_items table
CREATE TABLE inventory.order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES inventory.orders(id),
    product_id INTEGER REFERENCES inventory.products(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

-- Insert test data into products
INSERT INTO inventory.products (name, description, price, quantity) VALUES
('Laptop', 'High-performance laptop with 16GB RAM', 1299.99, 50),
('Mouse', 'Wireless ergonomic mouse', 29.99, 200),
('Keyboard', 'Mechanical gaming keyboard', 89.99, 100),
('Monitor', '27-inch 4K monitor', 399.99, 75),
('Headphones', 'Noise-cancelling wireless headphones', 199.99, 120),
('USB Hub', '7-port USB 3.0 hub', 39.99, 150),
('Webcam', '1080p HD webcam', 59.99, 80),
('Microphone', 'USB condenser microphone', 79.99, 60),
('Desk Lamp', 'LED desk lamp with adjustable brightness', 34.99, 90),
('Cable Organizer', 'Cable management system', 15.99, 200);

-- Insert test data into customers
INSERT INTO inventory.customers (first_name, last_name, email, phone) VALUES
('John', 'Doe', 'john.doe@email.com', '555-0101'),
('Jane', 'Smith', 'jane.smith@email.com', '555-0102'),
('Bob', 'Johnson', 'bob.johnson@email.com', '555-0103'),
('Alice', 'Williams', 'alice.williams@email.com', '555-0104'),
('Charlie', 'Brown', 'charlie.brown@email.com', '555-0105'),
('Diana', 'Davis', 'diana.davis@email.com', '555-0106'),
('Eve', 'Miller', 'eve.miller@email.com', '555-0107'),
('Frank', 'Wilson', 'frank.wilson@email.com', '555-0108'),
('Grace', 'Moore', 'grace.moore@email.com', '555-0109'),
('Henry', 'Taylor', 'henry.taylor@email.com', '555-0110');

-- Insert test data into orders
INSERT INTO inventory.orders (customer_id, status, total_amount) VALUES
(1, 'COMPLETED', 1329.98),
(2, 'PENDING', 89.99),
(3, 'SHIPPED', 459.98),
(4, 'COMPLETED', 239.98),
(5, 'PENDING', 1699.97),
(6, 'CANCELLED', 79.99),
(7, 'COMPLETED', 429.97),
(8, 'SHIPPED', 94.98),
(9, 'PENDING', 199.99),
(10, 'COMPLETED', 119.98);

-- Insert test data into order_items
INSERT INTO inventory.order_items (order_id, product_id, quantity, price) VALUES
(1, 1, 1, 1299.99),
(1, 2, 1, 29.99),
(2, 3, 1, 89.99),
(3, 4, 1, 399.99),
(3, 7, 1, 59.99),
(4, 5, 1, 199.99),
(4, 6, 1, 39.99),
(5, 1, 1, 1299.99),
(5, 4, 1, 399.99),
(6, 8, 1, 79.99),
(7, 9, 3, 34.99),
(7, 10, 2, 15.99),
(7, 2, 1, 29.99),
(7, 3, 3, 89.99),
(8, 8, 1, 79.99),
(8, 10, 1, 15.99),
(9, 5, 1, 199.99),
(10, 2, 2, 29.99),
(10, 7, 1, 59.99);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for products table
CREATE TRIGGER update_products_updated_at BEFORE UPDATE
    ON inventory.products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Grant necessary permissions for replication
ALTER TABLE inventory.products REPLICA IDENTITY FULL;
ALTER TABLE inventory.customers REPLICA IDENTITY FULL;
ALTER TABLE inventory.orders REPLICA IDENTITY FULL;
ALTER TABLE inventory.order_items REPLICA IDENTITY FULL;