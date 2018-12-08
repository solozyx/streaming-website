# mysql 5.5

CREATE TABLE users
(
    id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
    login_name varchar(64),
    pwd text NOT NULL
);
CREATE UNIQUE INDEX users_id_uindex ON users (id);
CREATE UNIQUE INDEX users_login_name_uindex ON users (login_name);
ALTER TABLE users COMMENT = '用户表';

CREATE TABLE video_info
(
    id varchar(64) PRIMARY KEY NOT NULL,
    author_id int,
    name text,
    display_ctime text,
    create_time timestamp DEFAULT NOW()
);
ALTER TABLE video_info COMMENT = '视频信息表';

CREATE TABLE comments
(
    id varchar(64) PRIMARY KEY NOT NULL,
    video_id varchar(64),
    author_id int,
    content text,
    time datetime
);
ALTER TABLE comments COMMENT = '视频评论表';

CREATE TABLE sessions
(
    session_id tinytext NOT NULL,
    TTL tinytext,
    login_name text
);
ALTER TABLE sessions COMMENT = '用户登录session表';

