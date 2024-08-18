PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE [requests] (
id INTEGER NOT NULL PRIMARY KEY,
time DATETIME NOT NULL
, headers TEXT NOT NULL, body TEXT NOT NULL, host TEXT NOT NULL, "method" TEXT NOT NULL, bin text not null);
CREATE UNIQUE INDEX bin_index on requests(bin);
COMMIT;
