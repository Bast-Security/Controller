CREATE TABLE Permissions (
    PermissionId INT NOT NULL AUTO INCREMENT,
    RoleName VARCHAR(16) NOT NULL,
    PRIMARY KEY (PermissionId),
    FOREIGN KEY (RoleName) REFERENCES Roles.Name
);

