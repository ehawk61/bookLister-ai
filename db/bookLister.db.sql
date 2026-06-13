BEGIN TRANSACTION;
DROP TABLE IF EXISTS "books";
CREATE TABLE IF NOT EXISTS "books" (
	"book_id"	INTEGER,
	"title"	TEXT,
	"category"	TEXT NOT NULL DEFAULT 'Fiction' CHECK("category" IN ('Non-Fiction', 'Fiction')),
	PRIMARY KEY("book_id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "authors_books";
CREATE TABLE IF NOT EXISTS "authors_books" (
	"authors_books_id"	INTEGER,
	"book_id"	INTEGER NOT NULL,
	"author_id"	INTEGER NOT NULL,
	PRIMARY KEY("authors_books_id" AUTOINCREMENT),
	FOREIGN KEY("author_id") REFERENCES "authors"("author_id"),
	FOREIGN KEY("book_id") REFERENCES "books"("book_id")
);
DROP TABLE IF EXISTS "books_read";
CREATE TABLE IF NOT EXISTS "books_read" (
	"books_read_id"	INTEGER,
	"book_id"	INTEGER NOT NULL,
	"year_read"	INTEGER NOT NULL,
	"readthrough_number"	INTEGER NOT NULL,
	FOREIGN KEY("book_id") REFERENCES "books"("book_id"),
	PRIMARY KEY("books_read_id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "authors";
CREATE TABLE IF NOT EXISTS "authors" (
	"author_id"	INTEGER,
	"prefix"	TEXT,
	"first_name"	TEXT,
	"middle_name"	TEXT,
	"last_name"	TEXT,
	"suffix"	TEXT,
	"organization"	TEXT,
	PRIMARY KEY("author_id" AUTOINCREMENT)
);
COMMIT;