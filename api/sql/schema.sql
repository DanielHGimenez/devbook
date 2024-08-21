create database if not exists devbook;
use devbook;

drop table if exists likes;
drop table if exists posts;
drop table if exists followers;
drop table if exists users;

create table users (
    id int auto_increment primary key,
    name varchar(50) not null,
    nick varchar(50) not null unique,
    email varchar(50) not null unique,
    password varchar(72) not null unique,
    created_at timestamp default current_timestamp()
) engine=innodb;

create table followers (
	followed_id int not null,
	foreign key (followed_id)
	references users(id)
	on delete cascade,

	follower_id int not null,
	foreign key (follower_id)
	references users(id)
	on delete cascade,

	primary key(followed_id, follower_id)
) engine=innodb;

create table posts (
	id int auto_increment primary key,
	title varchar(50) not null,
	content varchar(250) not null,
	likes int default 0,

	author_id int not null,
	foreign key (author_id)
	references users(id)
	on delete cascade,

	created_at timestamp default current_timestamp()
) engine=innodb;

create table likes (
	user_id int not null,
	foreign key (user_id)
	references users(id)
	on delete cascade,

	post_id int not null,
	foreign key (post_id)
	references posts(id)
	on delete cascade,

	primary key(user_id, post_id)
) engine=innodb;
