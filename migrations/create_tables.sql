
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
    id INTEGER AUTO_INCREMENT NOT NULL,
    role VARCHAR(32) NOT NULL,
    door VARCHAR(32) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (role) REFERENCES Roles (name),
    FOREIGN KEY (door) REFERENCES Doors (name)
);

CREATE TABLE AuthOption (
    id INTEGER AUTO_INCREMENT NOT NULL,
    door VARCHAR(32) NOT NULL
    PRIMARY KEY (id),
    FOREIGN KEY (door) REFERENCES Doors (name)
);

CREATE TABLE AuthType (
    option INTEGER NOT NULL,
    type INTEGER NOT NULL,
    PRIMARY KEY (option, type),
    FOREIGN KEY (option) REFERENCES AuthOption (id)
);

