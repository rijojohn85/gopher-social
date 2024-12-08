CREATE TABLE IF NOT EXISTS roles(
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  level INT NOT NULL,
  description VARCHAR(1000) NOT NULL
);

INSERT INTO roles (
  name, level, description
) VALUES ( 'user', 1, 'can make posts, delete and modify his posts and view all posts' );
INSERT INTO roles (
  name, level, description
) VALUES ( 'moderator', 2, 'can modify all users posts' );
INSERT INTO roles (
  name, level, description
) VALUES ( 'admin', 3, 'can modify and delete all posts' );
