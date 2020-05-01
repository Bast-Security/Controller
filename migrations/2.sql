
CREATE TABLE IF NOT EXISTS History (
	door INTEGER NOT NULL,
	time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	pin VARCHAR(32),
	card VARCHAR(32),
	PRIMARY KEY (door, time),
	FOREIGN KEY (door) REFERENCES Doors (id)
);
