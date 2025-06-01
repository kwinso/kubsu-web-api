CREATE TABLE
  IF NOT EXISTS languages (id SERIAL PRIMARY KEY, name VARCHAR(20) NOT NULL);

CREATE TABLE
  IF NOT EXISTS submissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    phone VARCHAR(500) NOT NULL,
    email VARCHAR(500) NOT NULL,
    birth_date VARCHAR(10) NOT NULL,
    bio TEXT NOT NULL,
    sex SMALLINT NOT NULL, -- Use 0/1 or consider BOOLEAN if appropriate
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username VARCHAR(12) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL
  );

CREATE TABLE
  IF NOT EXISTS submission_languages (
    submission_id INT NOT NULL,
    language_id INT NOT NULL,
    PRIMARY KEY (submission_id, language_id),
    FOREIGN KEY (submission_id) REFERENCES submissions (id),
    FOREIGN KEY (language_id) REFERENCES languages (id)
  );

INSERT INTO
  languages (name)
VALUES
  ('Pascal'),
  ('C'),
  ('C++'),
  ('JavaScript'),
  ('PHP'),
  ('Python'),
  ('Java'),
  ('Haskell'),
  ('Clojure'),
  ('Prolog'),
  ('Scala'),
  ('Go');

-- Insert an admin with a hashed password
-- Credentials: admin:admin
-- Hashed using bcrypt with cost factor 12
-- INSERT INTO admins (username, password) VALUES
-- (
--     'admin',
--     '$2y$12$XpzZch1FcMM9YXK2d9bygeWbwwhHcacdeJFaYGRDkS4DPAsq7SReW'
-- );