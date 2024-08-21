insert into users (name, nick, email, password)
values
    ('user 1', 'user_1', 'user1@gmail.com', '$2a$10$1XtWcaMz96Fhe9IqKhWFq.KBpZ5cYWV.sfZWRLfCRSZorYZwtEbgG'), -- pasword: password123
    ('user 2', 'user_2', 'user2@gmail.com', '$2a$10$1XtWcaMz96Fhe9IqKhWFq.KBpZ5cYWV.sfZWRLfCRSZorYZwtEbgG'),
    ('user 3', 'user_3', 'user3@gmail.com', '$2a$10$1XtWcaMz96Fhe9IqKhWFq.KBpZ5cYWV.sfZWRLfCRSZorYZwtEbgG');

insert into followers (user_id, follower_id)
values
    (1, 2),
    (1, 3),
    (3, 1);