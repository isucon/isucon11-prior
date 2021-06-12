DROP DATABASE IF EXISTS `isucon2021_prior`;
CREATE DATABASE IF NOT EXISTS `isucon2021_prior` DEFAULT CHARACTER SET utf8mb4;
CREATE USER IF NOT EXISTS `isucon`@`localhost` IDENTIFIED WITH mysql_native_password BY 'isucon';
GRANT ALL ON `isucon2021_prior`.* TO `isucon`@`localhost`;
