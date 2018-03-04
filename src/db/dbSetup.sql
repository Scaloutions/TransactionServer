CREATE TABLE IF NOT EXISTS users (
  user_id           VARCHAR(32) PRIMARY KEY,
  user_name         VARCHAR(20) UNIQUE NOT NULL,
  -- account_number    VARCHAR(32) UNIQUE NOT NULL,
  user_address      VARCHAR(100),
  user_email        VARCHAR(50)
);

/*
  TODO: we want to store dollars and cents separately
        and set the rigth constraints on numerical fields
*/
CREATE TABLE IF NOT EXISTS accounts (
  user_id        VARCHAR(32) PRIMARY KEY,
  balance           FLOAT(18,8),
  available_balance FLOAT(18,8)
);

CREATE TABLE IF NOT EXISTS stock (
  user_id        VARCHAR(32),
  symbol            VARCHAR(50),
  amount            FLOAT(18,8),
  PRIMARY KEY (user_id, symbol)
);