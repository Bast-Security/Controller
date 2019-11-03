CREATE TABLE Users (
    UserId INT NOT NULL AUTO INCREMENT 
    Name VARCHAR(32) NOT NULL,
    PubKey BINARY(32),
    PIN VARCHAR(4),
    CardNumber INT,
    LastAccess DATETIME,
    PRIMARY KEY (Id)
);

