-- MySQL dump 10.17  Distrib 10.3.22-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: bast
-- ------------------------------------------------------
-- Server version	10.3.22-MariaDB-0+deb10u1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `AuthOption`
--

DROP TABLE IF EXISTS `AuthOption`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `AuthOption` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `door` varchar(32) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `AuthOption`
--

LOCK TABLES `AuthOption` WRITE;
/*!40000 ALTER TABLE `AuthOption` DISABLE KEYS */;
/*!40000 ALTER TABLE `AuthOption` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CredentialType`
--

DROP TABLE IF EXISTS `CredentialType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CredentialType` (
  `authoptionid` int(11) NOT NULL,
  `type` int(11) NOT NULL,
  PRIMARY KEY (`authoptionid`,`type`),
  CONSTRAINT `CredentialType_ibfk_1` FOREIGN KEY (`authoptionid`) REFERENCES `AuthOption` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CredentialType`
--

LOCK TABLES `CredentialType` WRITE;
/*!40000 ALTER TABLE `CredentialType` DISABLE KEYS */;
/*!40000 ALTER TABLE `CredentialType` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Doors`
--

DROP TABLE IF EXISTS `Doors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Doors` (
  `name` varchar(32) NOT NULL,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Doors`
--

LOCK TABLES `Doors` WRITE;
/*!40000 ALTER TABLE `Doors` DISABLE KEYS */;
INSERT INTO `Doors` VALUES ('Main Entrance'),('Rear Entrance'),('Research Lab'),('Server Room'),('Side Entrance');
/*!40000 ALTER TABLE `Doors` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Permissions`
--

DROP TABLE IF EXISTS `Permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role` varchar(32) NOT NULL,
  `door` varchar(32) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `role` (`role`),
  KEY `door` (`door`),
  CONSTRAINT `Permissions_ibfk_1` FOREIGN KEY (`role`) REFERENCES `Roles` (`name`),
  CONSTRAINT `Permissions_ibfk_2` FOREIGN KEY (`door`) REFERENCES `Doors` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Permissions`
--

LOCK TABLES `Permissions` WRITE;
/*!40000 ALTER TABLE `Permissions` DISABLE KEYS */;
INSERT INTO `Permissions` VALUES (3,'Admin','Research Lab'),(4,'Admin','Server Room'),(14,'IT','Server Room'),(15,'Staff','Main Entrance'),(16,'Staff','Rear Entrance'),(17,'Staff','Side Entrance'),(18,'Guest','Main Entrance');
/*!40000 ALTER TABLE `Permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Roles`
--

DROP TABLE IF EXISTS `Roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Roles` (
  `name` varchar(32) NOT NULL,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Roles`
--

LOCK TABLES `Roles` WRITE;
/*!40000 ALTER TABLE `Roles` DISABLE KEYS */;
INSERT INTO `Roles` VALUES ('Admin'),('Developer'),('Guest'),('IT'),('Researcher'),('Staff');
/*!40000 ALTER TABLE `Roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `UserRole`
--

DROP TABLE IF EXISTS `UserRole`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `UserRole` (
  `userid` int(11) NOT NULL,
  `role` varchar(32) NOT NULL,
  PRIMARY KEY (`userid`,`role`),
  KEY `role` (`role`),
  CONSTRAINT `UserRole_ibfk_1` FOREIGN KEY (`userid`) REFERENCES `Users` (`id`),
  CONSTRAINT `UserRole_ibfk_2` FOREIGN KEY (`role`) REFERENCES `Roles` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `UserRole`
--

LOCK TABLES `UserRole` WRITE;
/*!40000 ALTER TABLE `UserRole` DISABLE KEYS */;
INSERT INTO `UserRole` VALUES (1,'Admin'),(2,'Developer'),(3,'IT'),(3,'Researcher'),(4,'Researcher'),(5,'Staff');
/*!40000 ALTER TABLE `UserRole` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Users`
--

DROP TABLE IF EXISTS `Users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `email` varchar(32) NOT NULL,
  `pin` varchar(32) DEFAULT NULL,
  `cardno` int(11) DEFAULT NULL,
  `lastaccess` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Users`
--

LOCK TABLES `Users` WRITE;
/*!40000 ALTER TABLE `Users` DISABLE KEYS */;
INSERT INTO `Users` VALUES (1,'Fabio','fabio@gmail.com','1234',1234,'2020-02-25 17:38:41'),(2,'Evan','evan@yahoo.com','4321',4321,'2020-02-25 17:38:41'),(3,'Lety','lety@outlook.com','5678',5678,'2020-02-25 17:38:41'),(4,'Kristen','kristen@aol.com','8765',8765,'2020-02-25 17:38:41'),(5,'Joe','joe@pm.com','1234',1234,'2020-02-25 17:41:50');
/*!40000 ALTER TABLE `Users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-02-25 10:35:50
