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
   accrual REAL DEFAULT 0,
   uploaded_at TIMESTAMP,
   CONSTRAINT fk_user
      FOREIGN KEY (user_id)
      REFERENCES users(user_id)
      ON DELETE CASCADE
);
