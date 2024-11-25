CREATE TABLE [bins] (
	bin_id INTEGER PRIMARY KEY,
	created_at DATETIME NOT NULL,
	owner TEXT
);
CREATE TABLE [requests] (
	id INTEGER PRIMARY KEY,
	timestamp DATETIME NOT NULL,
	headers TEXT NOT NULL,
	body TEXT NOT NULL,
	host TEXT NOT NULL,
  remoteAddr TEXT NOT NULL,
  requestUri TEXT NOT NULL,
	"method" TEXT NOT NULL,
	bin INTEGER NOT NULL,
	FOREIGN KEY (bin) REFERENCES bins(bin_id)	
);
