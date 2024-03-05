SET timezone = 'Asia/Jakarta';
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO roles (name)
VALUES
    ('admin'),
    ('computizen'),
    ('PUFA Computing'),
    ('PUMA IT'),
    ('PUMA IS'),
    ('PUMA VCD'),
    ('PUMA ID');
