USE ino;

CREATE TABLE tokens (
  id INT NOT NULL AUTO_INCREMENT,
  role VARCHAR(255) NOT NULL,
  role_id VARCHAR(255) NOT NULL,
  push_token VARCHAR(255) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE (role,role_id)
);