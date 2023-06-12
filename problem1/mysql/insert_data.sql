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