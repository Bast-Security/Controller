CREATE TABLE ActiveTimes (
    StartTime DATETIME NOT NULL,
    EndTime DATETIME NOT NULL,
    PermissionId INT NOT NULL,
    PRIMARY KEY (StartTime, EndTime),
    FOREIGN KEY (PermissionId) REFERENCES Permissions.PermissionId
);

