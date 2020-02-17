CREATE DATABASE tatltest;
USE tatltest;

CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `email` varchar(64) NOT NULL,
  `password` varchar(32) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
);

INSERT INTO user (email, password) VALUES("tatl@test.com", "e10adc3949ba59abbe56e057f20f883e");