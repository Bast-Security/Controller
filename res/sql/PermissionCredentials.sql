CREATE TABLE PermissionCredentials (
    PermissionId INT NOT NULL,
    CredentialType INT NOT NULL,
    PRIMARY KEY (PermissionId, CredentialType),
    FOREIGN KEY (PermissionId) REFERENCES Permissions.PermissionId
);
