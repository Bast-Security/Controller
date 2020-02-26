INSERT INTO `Doors` VALUES ('Main Entrance'),('Rear Entrance'),('Research Lab'),('Server Room'),('Side Entrance');

INSERT INTO `Roles` VALUES ('Admin'),('Developer'),('Guest'),('IT'),('Researcher'),('Staff');

INSERT INTO `Users` VALUES (1,'Fabio','fabio@gmail.com','1234',1234,'2020-02-25 17:38:41'),(2,'Evan','evan@yahoo.com','4321',4321,'2020-02-25 17:38:41'),(3,'Lety','lety@outlook.com','5678',5678,'2020-02-25 17:38:41'),(4,'Kristen','kristen@aol.com','8765',8765,'2020-02-25 17:38:41'),(5,'Joe','joe@pm.com','1234',1234,'2020-02-25 17:41:50');

INSERT INTO `UserRole` VALUES (1,'Admin'),(2,'Developer'),(3,'IT'),(3,'Researcher'),(4,'Researcher'),(5,'Staff');

INSERT INTO `Permissions` VALUES ('Admin','Research Lab'),('Admin','Server Room'),('Staff','Main Entrance'),('Staff','Rear Entrance'),('Staff','Side Entrance'),('Guest','Main Entrance');

-- AuthTypes 1 = PIN Only, 2 = Card Only, 3 = Card and PIN, -3 = Card or PIN
INSERT INTO AuthTypes VALUES ("Research Lab", -3), ("Server Room", 3), ("Main Entrance", 2), ("Side Entrance", 2), ("Rear Entrance", 2);
