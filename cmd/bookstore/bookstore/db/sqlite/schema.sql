PRAGMA strict = ON;
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

CREATE TABLE IF NOT EXISTS books (
	isbn TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	author TEXT NOT NULL,
	price INTEGER NOT NULL
);
