CREATE TABLE IF NOT EXISTS users (
  id       VARCHAR(100) PRIMARY KEY UNIQUE,
  name     VARCHAR(100),
  email    VARCHAR(100),
  password VARCHAR(100)
);
