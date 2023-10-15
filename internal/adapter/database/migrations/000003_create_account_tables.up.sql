CREATE TABLE IF NOT EXISTS withdrawals (
   withdraw_id SERIAL PRIMARY KEY,
   user_id INT NOT NULL,
   amount INT NOT NULL,
   order_number VARCHAR(100) UNIQUE NOT NULL,
   processed_at TIMESTAMP,
   CONSTRAINT fk_user
      FOREIGN KEY (user_id)
      REFERENCES users(user_id)
      ON DELETE CASCADE,
);

CREATE TABLE IF NOT EXISTS transactions (
   transaction_id SERIAL PRIMARY KEY,
   user_id INT UNIQUE NOT NULL,
   source_id INT NOT NULL,
   balance REAL,
   amount REAL,
   normal INT,
   updated_at TIMESTAMP DEFAULT NOW(),
   CONSTRAINT fk_user
      FOREIGN KEY (user_id)
      REFERENCES users(user_id)
      ON DELETE CASCADE,
   UNIQUE(source_id, normal)
);
