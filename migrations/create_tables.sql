CREATE TABLE IF NOT EXISTS Systems (
    id INTEGER NOT NULL AUTO_INCREMENT,
    name VARCHAR(32) NOT NULL,
    totpKey BINARY(32) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS Doors (
    id INTEGER NOT NULL AUTO_INCREMENT,
    name VARCHAR(32) NOT NULL,
    system INTEGER NOT NULL,
    keyX BLOB NOT NULL,
    keyY BLOB NOT NULL,
    challenge BINARY(16),
    PRIMARY KEY (id),
    FOREIGN KEY (system) REFERENCES Systems (id)
);

CREATE TABLE IF NOT EXISTS Users (
    id INTEGER NOT NULL AUTO_INCREMENT,
    system INTEGER NOT NULL,
    name VARCHAR(32) NOT NULL,
    email VARCHAR(32) NOT NULL,
    pin VARCHAR(32),
    cardno INTEGER,
    PRIMARY KEY(id, system),
    FOREIGN KEY (system) REFERENCES Systems (id)
);

CREATE TABLE IF NOT EXISTS Roles (
    name VARCHAR(32) NOT NULL,
    system INTEGER NOT NULL,
    PRIMARY KEY (name, system),
    FOREIGN KEY (system) REFERENCES Systems (id)
);

CREATE TABLE IF NOT EXISTS UserRole (
    system INTEGER NOT NULL,
    userid INTEGER NOT NULL,
    role VARCHAR(32) NOT NULL,
    PRIMARY KEY (system, userid, role),
    FOREIGN KEY (userid) REFERENCES Users(id),
    FOREIGN KEY (role) REFERENCES Roles(name),
    FOREIGN KEY (system) REFERENCES Systems (id)
);

CREATE TABLE IF NOT EXISTS Permissions (
    system INTEGER NOT NULL,
    role VARCHAR(32) NOT NULL,
    door INTEGER NOT NULL,
    PRIMARY KEY (role, door, system),
    FOREIGN KEY (system) REFERENCES Systems (id),
    FOREIGN KEY (role) REFERENCES Roles (name),
    FOREIGN KEY (door) REFERENCES Doors (id)
);

CREATE TABLE IF NOT EXISTS Admins (
    id INTEGER NOT NULL AUTO_INCREMENT,
    keyX BLOB NOT NULL,
    keyY BLOB NOT NULL,
    challenge BINARY(16),
    name VARCHAR(32),
    email VARCHAR(32),
    phone VARCHAR(32),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS AdminSystem (
    admin INTEGER NOT NULL,
    system INTEGER NOT NULL,
    PRIMARY KEY (admin, system),
    FOREIGN KEY (admin) REFERENCES Admins (id),
    FOREIGN KEY (system) REFERENCES Systems (id)
);

