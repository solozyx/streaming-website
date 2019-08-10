-- mysql 5.7.26 version

CREATE DATABASE `video_server` DEFAULT CHARACTER SET utf8;

CREATE TABLE `users`(
    `id` int NOT NULL AUTO_INCREMENT,
    `login_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
    `pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户密码',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_index_login_name` (`login_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';

CREATE TABLE `video_info`(
    `id` varchar(64) NOT NULL DEFAULT '' COMMENT '视频唯一标识',
    `author_id` int NOT NULL DEFAULT '0',
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '视频名称',
    `display_ctime` varchar(64) NOT NULL DEFAULT '' COMMENT '视频显示创建时间',
    `create_time` datetime DEFAULT NOW() ON UPDATE CURRENT_TIMESTAMP() COMMENT '视频创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='视频表';

CREATE TABLE `comments`(
    `id` varchar(64) NOT NULL,
    `video_id` varchar(64) NOT NULL DEFAULT '' COMMENT '视频唯一标识',
    `author_id` int NOT NULL DEFAULT '0' COMMENT '评论作者标识',
    `content` text,
    `time` datetime DEFAULT NOW() ON UPDATE CURRENT_TIMESTAMP() COMMENT '评论创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='视频评论表';

CREATE TABLE `sessions`(
    `session_id` varchar(64) NOT NULL DEFAULT '' COMMENT '会话session凭证',
    `TTL` varchar(64) NOT NULL DEFAULT '' COMMENT '过期时间',
    `login_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
    PRIMARY key (`session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户登录session表';
