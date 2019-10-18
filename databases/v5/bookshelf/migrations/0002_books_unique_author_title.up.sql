BEGIN;
ALTER TABLE IF EXISTS books DROP CONSTRAINT IF EXISTS books_unique_author_title;
ALTER TABLE IF EXISTS books ADD CONSTRAINT books_unique_author_title UNIQUE (author, title);
COMMIT;
