CREATE TABLE IF NOT EXISTS lists(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(32) NOT NULL
);

CREATE TABLE IF NOT EXISTS levels(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  list_id BIGINT NOT NULL REFERENCES lists(id),
  name VARCHAR(20) NOT NULL,
  description VARCHAR(255),
  showcase_link TEXT NOT NULL,
  placement INTEGER NOT NULL
    CHECK ( placement > 0 ),

  UNIQUE (list_id, placement)
);

CREATE TABLE IF NOT EXISTS users(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username VARCHAR(32) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(16) NOT NULL DEFAULT 'player'
    CHECK ( role IN ('player', 'moderator', 'owner') )
);

CREATE TABLE IF NOT EXISTS records(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  level_id BIGINT NOT NULL REFERENCES levels(id),
  status VARCHAR(16) NOT NULL DEFAULT 'pending'
    CHECK ( status IN ('pending', 'accepted', 'rejected') ),
  percentage INTEGER NOT NULL
    CHECK ( percentage BETWEEN 0 AND 100 ),
  video_link TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_permissions(
  user_id BIGINT NOT NULL REFERENCES users(id),
  list_id BIGINT NOT NULL REFERENCES lists(id),
  permissions SMALLINT NOT NULL DEFAULT 0
    CHECK ( permissions >= 0 ),

  PRIMARY KEY (user_id, list_id)
);

