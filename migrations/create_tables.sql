
CREATE TABLE Doors (
    name VARCHAR(32) NOT NULL,
    PRIMARY KEY (name)
);

CREATE TABLE Users (
    id INTEGER NOT NULL AUTO_INCREMENT,
    name VARCHAR(32) NOT NULL,
    email VARCHAR(32) NOT NULL,
    pin VARCHAR(32),
    cardno INTEGER,
    lastaccess TIMESTAMP,
    PRIMARY KEY(id)
);

CREATE TABLE Roles (
    name VARCHAR(32) NOT NULL,
    PRIMARY KEY (name)
);

CREATE TABLE UserRole (
    userid INTEGER NOT NULL,
    role VARCHAR(32) NOT NULL,
    PRIMARY KEY (userid, role),
    FOREIGN KEY (userid) REFERENCES Users(id),
    FOREIGN KEY (role) REFERENCES Roles(name)
);

CREATE TABLE Permissions (
    role VARCHAR(32) NOT NULL,
    door VARCHAR(32) NOT NULL,
    PRIMARY KEY (role, door),
    FOREIGN KEY (role) REFERENCES Roles (name),
    FOREIGN KEY (door) REFERENCES Doors (name)
);

CREATE TABLE AuthTypes (
    door VARCHAR(32) NOT NULL,
    authType INTEGER NOT NULL,
    PRIMARY KEY (door, authtype),
    FOREIGN KEY (door) REFERENCES Doors (name)
);

CREATE TABLE Settings (
    name VARCHAR(32) NOT NULL,
    value VARCHAR(32) NOT NULL,
    PRIMARY KEY (name)
);

CREATE TABLE Admins (
    id INTEGER AUTO_INCREMENT NOT NULL,
    pubKey BINARY(32) NOT NULL,
    challenge BINARY(16),
    PRIMARY KEY (id)
);


