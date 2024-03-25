CREATE TABLE users(
	id UUID PRIMARY KEY,
	username TEXT NOT NULL,
	email TEXT NOT NULL,
	password TEXT NOT NULL
);

CREATE UNIQUE INDEX users_username ON users(username);
