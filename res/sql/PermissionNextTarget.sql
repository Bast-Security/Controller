CREATE TABLE PermissionNextTarget (
    DoorName VARCHAR(32) NOT NULL,
    PermissionId INT NOT NULL,
    PRIMARY KEY (Door, PermissionId),
    FOREIGN KEY (Door) REFERENCES Doors.Name,
    FOREIGN KEY (PermissionId) REFERENCES Permissions.PermissionId
);

