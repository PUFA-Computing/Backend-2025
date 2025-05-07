SET timezone = 'Asia/Jakarta';
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO roles (name)
VALUES
    ('admin'),
    ('computizen'),
    ('PUFA Computer Science'),
    ('PUMA IT'),
    ('PUMA IS'),
    ('guest');