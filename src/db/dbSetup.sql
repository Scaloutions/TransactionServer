CREATE DATABASE DAYTRADING;

CREATE TABLE IF NOT EXISTS DAYTRADING.users (
  user_id           VARCHAR(32) PRIMARY KEY,
  user_name         VARCHAR(20) NOT NULL,
  -- account_number    VARCHAR(32) UNIQUE NOT NULL,
  user_address      VARCHAR(100),
  user_email        VARCHAR(50)
);

/*
  TODO: we want to store dollars and cents separately
        and set the rigth constraints on numerical fields
*/
CREATE TABLE IF NOT EXISTS DAYTRADING.accounts (
  user_id        VARCHAR(32) PRIMARY KEY,
  balance           FLOAT(18,8),
  available_balance FLOAT(18,8)
);

CREATE TABLE IF NOT EXISTS DAYTRADING.stock (
  user_id        VARCHAR(32),
  symbol            VARCHAR(10),
  amount            FLOAT(18,8),
  available_amount  FLOAT(18,8),
  PRIMARY KEY (user_id, symbol)
);

CREATE TABLE IF NOT EXISTS DAYTRADING.buy_triggers (
  user_id         VARCHAR(32),
  stock           VARCHAR(10),
  stock_amount    FLOAT(18,8),
  --transaction_num     INT(),
  PRIMARY KEY (user_id, stock),
);

--CREATE INDEX buy_trigger ON DAYTRADING.buy_triggers (user_id, stock);

CREATE TABLE IF NOT EXISTS DAYTRADING.sell_triggers (
  user_id         VARCHAR(32),
  stock           VARCHAR(10),
  stock_amount    FLOAT(18,8),
  --transaction_num     INT(),
  PRIMARY KEY (user_id, stock),
);

--CREATE INDEX sell_trigger ON DAYTRADING.sell_triggers (user_id, stock);

-- TODO: should we track this with a 60sec timestamp?
CREATE TABLE IF NOT EXISTS DAYTRADING.sell (
  user_id         VARCHAR(32),
  stock           VARCHAR(10),
  stock_amount    FLOAT(18,8),
  money_amount    FLOAT(18,8),
  transaction_num     INT
);

CREATE INDEX sell_cmd ON DAYTRADING.sell (user_id, transaction_num);

CREATE TABLE IF NOT EXISTS DAYTRADING.buy (
  user_id         VARCHAR(32),
  stock           VARCHAR(10),
  stock_amount    FLOAT(18,8),
  money_amount    FLOAT(18,8),
  transaction_num     INT
);

CREATE INDEX buy_cmd ON DAYTRADING.buy (user_id, transaction_num);