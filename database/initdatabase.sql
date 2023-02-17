CREATE TABLE `users` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `full_name` varchar(50),
  `username` varchar(16) NOT NULL,
  `password` varchar(100) NOT NULL,
  `authority_type` varchar(20) DEFAULT "customer",
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `authorities` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `authority_type` varchar(20),
  `permissions` varchar(100),
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX `users_index_0` ON `users` (`usersname`);

CREATE UNIQUE INDEX `authorities_index_1` ON `authorities` (`authority_type`);

ALTER TABLE `users` ADD FOREIGN KEY (`authority_type`) REFERENCES `authorities` (`authority_type`);

insert INTO authorities(authority_type, permissions) value ("admin", "do anything");
insert INTO authorities(authority_type, permissions) value ("customer", "just buy stuff");
insert INTO authorities(authority_type, permissions) value ("employee", "restock goods");

insert INTO users(full_name, username, password, authority_type) value ("admin", "admin", "81dc9bdb52d04dc20036dbd8313ed055", "admin");
insert INTO users(full_name, username, password, authority_type) value ("customer", "customer", "81dc9bdb52d04dc20036dbd8313ed055", "customer");
insert INTO users(full_name, username, password, authority_type) value ("employee", "employee", "81dc9bdb52d04dc20036dbd8313ed055", "employee");