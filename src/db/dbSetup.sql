CREATE TABLE IF NOT EXISTS users (
  user_id           VARCHAR(32) PRIMARY KEY,
  user_name         VARCHAR(20) UNIQUE NOT NULL,
  account_number    VARCHAR(32) UNIQUE NOT NULL,
  user_address      VARCHAR(100),
  user_email        VARCHAR(50)
);