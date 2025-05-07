CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
    );

INSERT INTO organizations (name)
VALUES
    ('PUFA Computer Science'),
    ('PUMA IT'),
    ('PUMA IS');
