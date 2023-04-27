CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `name` varchar(64) DEFAULT '' NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT unique_user_id UNIQUE (user_id)
);
-- user1 user2
CREATE TABLE `friend_link` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user1_id` int(11) NOT NULL,
  `user2_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
);
-- user1 user2 block
CREATE TABLE `block_list` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user1_id` int(11) NOT NULL,
  `user2_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO users (user_id, name) VALUES (1, 'Alice');
INSERT INTO users (user_id, name) VALUES (2, 'Bob');
INSERT INTO users (user_id, name) VALUES (3, 'Charlie');
INSERT INTO users (user_id, name) VALUES (4, 'David');
INSERT INTO users (user_id, name) VALUES (5, 'Eve');
INSERT INTO users (user_id, name) VALUES (6, 'Kazuki');
INSERT INTO users (user_id, name) VALUES (7, 'NOBODY');

INSERT INTO friend_link (user1_id, user2_id) VALUES (1, 2);
INSERT INTO friend_link (user1_id, user2_id) VALUES (1, 3);
INSERT INTO friend_link (user1_id, user2_id) VALUES (2, 3);
INSERT INTO friend_link (user1_id, user2_id) VALUES (4, 5);
INSERT INTO friend_link (user1_id, user2_id) VALUES (5, 6);
INSERT INTO friend_link (user1_id, user2_id) VALUES (1, 6);


INSERT INTO block_list (user1_id, user2_id) VALUES (2, 3);
INSERT INTO block_list (user1_id, user2_id) VALUES (3, 1);
