CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,     
    last_name VARCHAR(100) NOT NULL,     
    email VARCHAR(255) NOT NULL UNIQUE,  
    city VARCHAR(100) NOT NULL,          
    password VARCHAR(255) NOT NULL        
);