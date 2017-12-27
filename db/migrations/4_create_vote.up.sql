CREATE TABLE IF NOT EXISTS vote (
  id         VARCHAR(100) PRIMARY KEY UNIQUE,
  created_at VARCHAR(100),
  user_id    VARCHAR(100) REFERENCES users (id),
  link_id    VARCHAR(100) REFERENCES link (id)
);
