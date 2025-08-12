-- create a table for users
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- add data for users
INSERT INTO users (name, email, password) VALUES
('John Doe', 'johndoe@example.com', 'password123'),
('Jane Smith', 'janesmith@example.com', 'password456'),
('Alice Johnson', 'alicejohnson@example.com', 'password789');

ALTER TABLE users ADD COLUMN password_changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP AFTER password;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE users
    ALTER COLUMN name TYPE VARCHAR(100),
    ALTER COLUMN password TYPE VARCHAR(255),
    ALTER COLUMN email TYPE VARCHAR(150),
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at::timestamp,
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at::timestamp;

ALTER TABLE users
    ALTER COLUMN email SET NOT NULL;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
    ADD COLUMN IF NOT EXISTS role VARCHAR(20),
    ADD COLUMN IF NOT EXISTS avatar VARCHAR(255),
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;

ALTER TABLE users
    ADD CONSTRAINT users_email_unique UNIQUE (email);

ALTER TABLE users
    ADD CONSTRAINT users_role_check CHECK (role IN ('rider', 'driver', 'admin'));

ALTER TABLE users
    ALTER COLUMN password TYPE TEXT;

-- Ensure column exists and has appropriate size (optional)
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS phone VARCHAR(20);

-- Add unique constraint on phone
ALTER TABLE users
    ADD CONSTRAINT users_phone_unique UNIQUE (phone);
