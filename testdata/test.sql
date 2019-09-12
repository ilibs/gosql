# Dump of table moments
# ------------------------------------------------------------

DROP TABLE IF EXISTS `moments`;

CREATE TABLE `moments` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '成员ID',
  `content` text NOT NULL COMMENT '日记内容',
  `comment_total` int(11) NOT NULL DEFAULT '0' COMMENT '评论总数',
  `like_total` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '1 正常 2删除',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `moments` WRITE;
/*!40000 ALTER TABLE `moments` DISABLE KEYS */;

INSERT INTO `moments` (`id`, `user_id`, `content`, `comment_total`, `like_total`, `status`, `created_at`, `updated_at`)
VALUES
	(1,5,'sdfsdfsdfsdfsdf',0,0,1,'2018-11-28 14:04:02','2018-11-28 14:04:02'),
	(2,5,'sdfsdfsdfsdfsdf',0,0,1,'2018-11-28 17:14:23','2018-11-28 17:14:23'),
	(3,6,'123123123',0,0,1,'2018-11-28 17:19:38','2018-11-28 17:19:38'),
	(4,5,'13212312313',0,0,1,'2018-11-28 17:22:25','2018-11-28 17:22:25'),
	(5,5,'123123123123',0,0,1,'2018-11-28 17:24:21','2018-11-28 17:24:21'),
	(6,6,'131231232345tasvdf',0,0,1,'2018-11-28 17:24:27','2018-11-28 17:24:27'),
	(7,5,'1231231231231231',0,0,1,'2018-11-28 18:07:48','2018-11-28 18:07:48'),
	(8,5,'1231231231231231',0,0,1,'2018-11-28 18:09:20','2018-11-28 18:09:20'),
	(9,6,'1231231231231231',0,0,1,'2018-11-28 18:11:19','2018-11-28 18:11:19'),
	(10,5,'1231231231231231',0,0,1,'2018-11-28 18:13:52','2018-11-28 18:13:52'),
	(11,5,'1231231231231231',0,0,1,'2018-11-28 18:15:02','2018-11-28 18:15:02'),
	(12,5,'1231231231231231',0,0,1,'2018-11-28 18:15:13','2018-11-28 18:15:13'),
	(13,5,'1231231231231231',0,0,1,'2018-11-28 18:15:39','2018-11-28 18:15:39'),
	(14,6,'开开信息想你',0,0,1,'2018-11-28 18:31:37','2018-11-28 18:31:37'),
	(15,5,'网友们已经开始争相给宝宝取名字了',0,0,1,'2018-11-28 18:34:45','2018-11-28 18:34:45'),
	(16,6,' B2B事业部的对外报价显示',0,0,1,'2018-11-28 18:35:24','2018-11-28 18:35:24');

/*!40000 ALTER TABLE `moments` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table photos
# ------------------------------------------------------------

DROP TABLE IF EXISTS `photos`;

CREATE TABLE `photos` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '照片路径',
  `moment_id` int(11) NOT NULL COMMENT '日记ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `photos` WRITE;
/*!40000 ALTER TABLE `photos` DISABLE KEYS */;

INSERT INTO `photos` (`id`, `url`, `moment_id`, `created_at`, `updated_at`)
VALUES
	(1,'https://static.fifsky.com/kids/upload/20181128/5febe3b6e23623168cb70eac39d26412.png!blog',1,'2018-11-28 18:15:39','2018-11-28 18:15:39'),
	(2,'https://static.fifsky.com/kids/upload/20181128/9c60f42f07d7a0e13293c91fc5740c9d.png!blog',10,'2018-11-28 18:15:39','2018-11-28 18:15:39'),
	(3,'https://static.fifsky.com/kids/upload/20181128/458762098fb20128996c9cb21309aa9a.png!blog',1,'2018-11-28 18:31:37','2018-11-28 18:31:37'),
	(4,'https://static.fifsky.com/kids/upload/20181128/5b90c5af1bc35375a08cbc990ed662d1.png!blog',14,'2018-11-28 18:31:37','2018-11-28 18:31:37'),
	(5,'https://static.fifsky.com/kids/upload/20181128/db190be4184774d88abb31521123b14c.png!blog',14,'2018-11-28 18:31:37','2018-11-28 18:31:37'),
	(6,'https://static.fifsky.com/kids/upload/20181128/e1bd15706a79edd1f92f54538622600e.png!blog',14,'2018-11-28 18:31:37','2018-11-28 18:31:37'),
	(7,'https://static.fifsky.com/kids/upload/20181128/6bf495726054fa12ae7e6f5d0d4560a4.png!blog',14,'2018-11-28 18:31:37','2018-11-28 18:31:37'),
	(8,'https://static.fifsky.com/kids/upload/20181128/e1bd15706a79edd1f92f54538622600e.png!blog',15,'2018-11-28 18:34:45','2018-11-28 18:34:45'),
	(9,'https://static.fifsky.com/kids/upload/20181128/c6cc28b912f805b6ef402603e0f67852.png!blog',9,'2018-11-28 18:34:45','2018-11-28 18:34:45'),
	(10,'https://static.fifsky.com/kids/upload/20181128/a63a0798d098272a39e76c88f39f2f29.png!blog',16,'2018-11-28 18:35:24','2018-11-28 18:35:24'),
	(11,'https://static.fifsky.com/kids/upload/20181128/381de1930d970183ab083fe08e2677ac.png!blog',16,'2018-11-28 18:35:24','2018-11-28 18:35:24');

/*!40000 ALTER TABLE `photos` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table moment_users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `moment_users`;

CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '',
  `status` int(11) NOT NULL,
  `success_time` datetime,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `moment_users` DISABLE KEYS */;

INSERT INTO `users` (`id`,`name`, `status`, `created_at`, `updated_at`)
VALUES
	(5,'豆爸&玥爸',1,'2018-11-28 10:29:55','2018-11-28 10:29:55'),
	(6,'呵呵',1,'2018-11-28 10:29:55','2018-11-28 10:29:55');

/*!40000 ALTER TABLE `moment_users` ENABLE KEYS */;
UNLOCK TABLES;
