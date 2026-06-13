DROP TABLE IF EXISTS "authors_books";
DROP TABLE IF EXISTS "books_read";
DROP TABLE IF EXISTS "books";
DROP TABLE IF EXISTS "authors";

CREATE TABLE IF NOT EXISTS "authors" (
	"author_id"	INTEGER GENERATED ALWAYS AS IDENTITY,
	"prefix"	TEXT,
	"first_name"	TEXT,
	"middle_name"	TEXT,
	"last_name"	TEXT,
	"suffix"	TEXT,
	"organization"	TEXT,
	PRIMARY KEY("author_id")
);

CREATE TABLE IF NOT EXISTS "books" (
	"book_id"	INTEGER GENERATED ALWAYS AS IDENTITY,
	"title"	TEXT,
	"category"	TEXT NOT NULL DEFAULT 'Fiction' CHECK("category" IN ('Non-Fiction', 'Fiction')),
	PRIMARY KEY("book_id")
);

CREATE TABLE IF NOT EXISTS "books_read" (
	"books_read_id"	INTEGER GENERATED ALWAYS AS IDENTITY,
	"book_id"	INTEGER NOT NULL,
	"year_read"	INTEGER NOT NULL,
	"readthrough_number"	INTEGER NOT NULL,
	PRIMARY KEY("books_read_id"),
	FOREIGN KEY("book_id") REFERENCES "books"("book_id")
);

CREATE TABLE IF NOT EXISTS "authors_books" (
	"authors_books_id"	INTEGER GENERATED ALWAYS AS IDENTITY,
	"book_id"	INTEGER NOT NULL,
	"author_id"	INTEGER NOT NULL,
	PRIMARY KEY("authors_books_id"),
	FOREIGN KEY("author_id") REFERENCES "authors"("author_id"),
	FOREIGN KEY("book_id") REFERENCES "books"("book_id")
);
