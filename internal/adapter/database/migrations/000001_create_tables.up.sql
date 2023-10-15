CREATE TABLE IF NOT EXISTS users (
   user_id serial PRIMARY KEY,
   login VARCHAR (100) UNIQUE NOT NULL,
   password VARCHAR (100) NOT NULL
);

CREATE TYPE order_status AS ENUM (
   'NEW',
   'PROCESSING',
   'INVALID',
   'PROCESSED'
);

CREATE TABLE IF NOT EXISTS orders (
   order_id SERIAL PRIMARY KEY,
   user_id INT,
   number VARCHAR(100) UNIQUE NOT NULL,
   status order_status NOT NULL,
   accrual REAL NULL,
   uploaded_at TIMESTAMP DEFAULT NOW(),
   CONSTRAINT fk_user
      FOREIGN KEY (user_id)
      REFERENCES users(user_id)
      ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS withdrawals (
   withdraw_id SERIAL PRIMARY KEY,
   user_id INT NOT NULL,
   amount REAL NOT NULL,
   order_number VARCHAR(100) UNIQUE NOT NULL,
   processed_at TIMESTAMP,
   CONSTRAINT fk_user
      FOREIGN KEY (user_id)
      REFERENCES users(user_id)
      ON DELETE CASCADE
);

CREATE TYPE transaction_type AS ENUM (
   'ACCRUAL',
   'WITHDRAW'
);

CREATE TABLE IF NOT EXISTS transactions (
   transaction_id SERIAL PRIMARY KEY,
   user_id INT NOT NULL,
   source_id INT NOT NULL,
   source_type transaction_type,
   balance REAL NOT NULL,
   created_at TIMESTAMP DEFAULT NOW(),
   CONSTRAINT fk_user
      FOREIGN KEY (user_id)
      REFERENCES users(user_id)
      ON DELETE CASCADE,
   UNIQUE(source_id, source_type)
);
