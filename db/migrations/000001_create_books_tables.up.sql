-- 書籍データベース用のテーブルをセットアップする -- :setup_table_books
CREATE TABLE books (
  isbn varchar(20) PRIMARY KEY,
  title varchar(256) NOT NULL, 
  author varchar(256) NOT NULL, 
  owned smallint DEFAULT 0, 
  available smallint DEFAULT 0
);