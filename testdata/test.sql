# Dump of table cates
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cates`;

CREATE TABLE `cates` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `desc` varchar(255) NOT NULL DEFAULT '',
  `domain` varchar(100) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `un_domain` (`domain`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

LOCK TABLES `cates` WRITE;
/*!40000 ALTER TABLE `cates` DISABLE KEYS */;

INSERT INTO `cates` (`id`, `name`, `desc`, `domain`, `created_at`, `updated_at`)
VALUES
	(1,'默认分类','默认分类','default','2017-08-18 15:21:56','2017-08-18 15:21:56'),
	(4,'技术笔记','技术笔记,PHP,REDIS,LINUX,MYSQL,GO','notes','2017-08-18 15:21:56','2017-08-18 15:21:56');

/*!40000 ALTER TABLE `cates` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table links
# ------------------------------------------------------------

DROP TABLE IF EXISTS `links`;

CREATE TABLE `links` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `url` varchar(200) NOT NULL DEFAULT '',
  `desc` varchar(255) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

LOCK TABLES `links` WRITE;
/*!40000 ALTER TABLE `links` DISABLE KEYS */;

INSERT INTO `links` (`id`, `name`, `url`, `desc`, `created_at`)
VALUES
	(1,'fifsky','http://fifsky.com','fifsky','2017-08-18 15:21:56');

/*!40000 ALTER TABLE `links` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table moods
# ------------------------------------------------------------

DROP TABLE IF EXISTS `moods`;

CREATE TABLE `moods` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(255) NOT NULL DEFAULT '',
  `user_id` int(10) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

LOCK TABLES `moods` WRITE;
/*!40000 ALTER TABLE `moods` DISABLE KEYS */;

INSERT INTO `moods` (`id`, `content`, `user_id`, `created_at`)
VALUES
	(1,'Hi,fifsky!',1,'2017-08-18 15:21:56');

/*!40000 ALTER TABLE `moods` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table options
# ------------------------------------------------------------

DROP TABLE IF EXISTS `options`;

CREATE TABLE `options` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `option_key` varchar(100) NOT NULL DEFAULT '',
  `option_value` varchar(200) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `option_name` (`option_key`) USING BTREE
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

LOCK TABLES `options` WRITE;
/*!40000 ALTER TABLE `options` DISABLE KEYS */;

INSERT INTO `options` (`id`, `option_key`, `option_value`)
VALUES
	(1,'site_name','無處告別'),
	(2,'site_desc','回首往事，珍重眼前人'),
	(3,'site_keyword','fifsky,rita,生活,博客,豆豆'),
	(4,'post_num','10');

/*!40000 ALTER TABLE `options` ENABLE KEYS */;
UNLOCK TABLES;