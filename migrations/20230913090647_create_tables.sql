-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id uuid DEFAULT gen_random_uuid() NOT NULL CONSTRAINT pk_user PRIMARY KEY,
  first_name VARCHAR(50) NOT NULL CHECK (first_name <> ''),
  last_name VARCHAR(100) NOT NULL CHECK (last_name <> ''),
  email VARCHAR(250) NOT NULL CHECK (email <> ''),
  password VARCHAR(255) NOT NULL CHECK (octet_length(password) <> 0),
  phone VARCHAR(50) NOT NULL CHECK (phone <> ''),
  profile_picture_url VARCHAR(255),
  city VARCHAR(50),
  birthday VARCHAR(50),
  age smallint CHECK (age >= 0 AND age <= 200),
  gender VARCHAR(50) NOT NULL CHECK (gender <> ''),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  last_login_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  role VARCHAR(50) NOT NULL,
  CONSTRAINT unique_email_phone UNIQUE (email, phone)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- +goose StatementEnd