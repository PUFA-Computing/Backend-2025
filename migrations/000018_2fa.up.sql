ALTER TABLE IF EXISTS users
ADD COLUMN IF NOT EXISTS "twofa_enabled" BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS "twofa_image" TEXT,
ADD COLUMN IF NOT EXISTS "twofa_secret" TEXT;


