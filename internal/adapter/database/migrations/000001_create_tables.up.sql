-- пользователи
CREATE TABLE IF NOT EXISTS users (
   user_id serial PRIMARY KEY,
   login VARCHAR (100) UNIQUE NOT NULL,
   password VARCHAR (100) NOT NULL
);

-- перечисление возможных статусов заказа
CREATE TYPE order_status AS ENUM (
   'NEW',
   'PROCESSING',
   'INVALID',
   'PROCESSED'
);

-- заказы
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

-- списания
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

-- перечисление возможных источников изменения транзакции
CREATE TYPE transaction_type AS ENUM (
   'ACCRUAL',
   'WITHDRAW'
);

-- транзакции
-- Источником транзакции может быть либо начисление от заказа, либо списание.
-- Транзакции пользователя должны применятся последовательно.
CREATE TABLE IF NOT EXISTS transactions (
   transaction_id SERIAL PRIMARY KEY,
   user_id INT NOT NULL,
   -- последовательный номер транзакции в рамках одного пользователя
   user_transaction_seq INT NOT NULL,
   -- источник транзакции
   source_id INT NOT NULL,
   -- тип источника транзакции
   source_type transaction_type,
   -- баланс в результате проведения транзакции
   balance REAL NOT NULL,
   created_at TIMESTAMP DEFAULT NOW(),
   CONSTRAINT fk_user
      FOREIGN KEY (user_id)
      REFERENCES users(user_id)
      ON DELETE CASCADE,
   -- одному источнику соответствует одна транзакция
   -- исключает двойное списание/начисление
   UNIQUE(source_id, source_type),
   -- транзакции пользователя могут применяться только последовательно
   -- позволяет избежать отрицательного баланса
   UNIQUE(user_id, user_transaction_seq)
);
