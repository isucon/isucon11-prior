CREATE DATABASE IF NOT EXISTS `isumark` DEFAULT CHARACTER SET utf8mb4;
CREATE USER IF NOT EXISTS `isucon`@`localhost` IDENTIFIED WITH mysql_native_password BY 'isucon';
GRANT ALL ON `isumark`.* TO `isucon`@`localhost`;
