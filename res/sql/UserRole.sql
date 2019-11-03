CREATE TABLE UserRole (
    UserId INT NOT NULL,
    Role INT NOT NULL,
    PRIMARY KEY (UserId, Role),
    FOREIGN KEY (UserId) REFERENCES Users.UserId,
    FOREIGN KEY (Role) REFERENCES Roles.Role
);

