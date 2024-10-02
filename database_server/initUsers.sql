-- Active: 1714235039086@@127.0.0.1@5432
-- Elimina la base de datos
DROP DATABASE IF EXISTS logs_db;
-- Crear la base de datos si no existe
CREATE DATABASE logs_db;

-- Conectar a la base de datos recién creada
\c logs_db;

-- Crear la tabla de usuarios si no existe
CREATE TABLE IF NOT EXISTS logs (
    id SERIAL PRIMARY KEY,
    app_name VARCHAR(255) NOT NULL,
    log_type VARCHAR(50) NOT NULL,
    module VARCHAR(255) NOT NULL,
    log_date_time TIMESTAMP NOT NULL,
    summary VARCHAR(255) NOT NULL,
    description TEXT NOT NULL
);

-- Insertar datos de prueba en la tabla logs
INSERT INTO logs (app_name, log_type, module, log_date_time, summary, description) VALUES
('USERS-API', 'ERROR', 'ADD-USER', '2024-04-01T10:00:00-05:00', 'Failed to create a user', 'No mandatory information was entered'),
('USERS-API', 'ERROR', 'ADD-USER', '2024-04-02T10:05:00-05:00', 'Failed to create a user', 'An email in use was entered'),
('USERS-API', 'INFO', 'ADD-USER', '2024-04-03T10:10:00-05:00', 'Create a user', 'A user has been created successfully'),
('USERS-API', 'ERROR', 'UPDATE-USER', '2024-04-04T10:15:00-05:00', 'Failed to update user', 'Invalid token provided'),
('USERS-API', 'ERROR', 'UPDATE-USER', '2024-04-05T10:20:00-05:00', 'Failed to update user', 'Email to update is already in use'),
('USERS-API', 'INFO', 'UPDATE-USER', '2024-04-06T10:25:00-05:00', 'User updated successfully', 'A user has been updated successfully'),
('USERS-API', 'ERROR', 'DELETE-USER', '2024-04-07T10:30:00-05:00', 'Failed to delete user', 'Email not found in the database'),
('USERS-API', 'ERROR', 'DELETE-USER', '2024-04-08T10:35:00-05:00', 'Failed to delete user', 'Invalid token provided'),
('USERS-API', 'INFO', 'DELETE-USER', '2024-04-09T10:40:00-05:00', 'User deleted successfully', 'User was deleted successfully'),
('USERS-API', 'ERROR', 'GET-USER-BY-ID', '2024-04-10T10:45:00-05:00', 'Failed to get user by ID', 'Invalid token provided'),
('USERS-API', 'ERROR', 'GET-USER-BY-ID', '2024-04-11T10:50:00-05:00', 'Failed to get user by ID', 'Failed to retrieve user ID from database'),
('USERS-API', 'INFO', 'GET-USER-BY-ID', '2024-04-12T10:55:00-05:00', 'User retrieved successfully', 'User ID retrieved successfully'),
('USERS-API', 'ERROR', 'GET-ALL-USERS', '2024-04-13T11:00:00-05:00', 'Failed to get all users', 'Invalid token provided'),
('USERS-API', 'INFO', 'GET-ALL-USERS', '2024-04-14T11:05:00-05:00', 'All users retrieved successfully', 'All users were retrieved successfully'),
('USERS-API', 'ERROR', 'LOGIN', '2024-04-15T11:10:00-05:00', 'Failed to log in', 'Method not allowed'),
('USERS-API', 'ERROR', 'LOGIN', '2024-04-16T11:15:00-05:00', 'Failed to log in', 'Missing required email or password'),
('USERS-API', 'ERROR', 'LOGIN', '2024-04-17T11:20:00-05:00', 'Failed to log in', 'User not found in the database'),
('USERS-API', 'INFO', 'LOGIN', '2024-04-17T11:25:00-05:00', 'Login successful', 'User logged in successfully');


--  psql -U admin -d logs_db -f init.sql