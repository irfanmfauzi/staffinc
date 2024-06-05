CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	email VARCHAR(20) UNIQUE,
	password VARCHAR(255) NOT NULL,
	role varchar(20) NOT NULL
);


CREATE TABLE IF NOT EXISTS generator_links (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT,
	code VARCHAR(255) NOT NULL,
	expired_at timestamp NOT NULL,

	FOREIGN KEY (user_ud) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_code ON generator_links(code);
CREATE INDEX idx_user_id ON generator_links(user_id);

