package gosql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/ilibs/gosql/v2/internal/example/models"
)

var (
	createSchemas = map[string]string{
		"moments": `
CREATE TABLE moments (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  user_id int(11) NOT NULL COMMENT '成员ID',
  content text NOT NULL COMMENT '日记内容',
  comment_total int(11) NOT NULL DEFAULT '0' COMMENT '评论总数',
  like_total int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
  status int(11) NOT NULL DEFAULT '1' COMMENT '1 正常 2删除',
  created_at datetime NOT NULL COMMENT '创建时间',
  updated_at datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		"users": `
CREATE TABLE users (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50) NOT NULL DEFAULT '',
  status int(11) NOT NULL,
  success_time datetime,
  created_at datetime NOT NULL,
  updated_at datetime NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`,
		"photos": `
CREATE TABLE photos (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  url varchar(255) NOT NULL DEFAULT '' COMMENT '照片路径',
  moment_id int(11) NOT NULL COMMENT '日记ID',
  created_at datetime NOT NULL COMMENT '创建时间',
  updated_at datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`,
	}

	datas = map[string]string{
		"users": `
INSERT INTO users (id,name, status, created_at, updated_at) VALUES
	(5,'豆爸&玥爸',1,'2018-11-28 10:29:55','2018-11-28 10:29:55'),
	(6,'呵呵',1,'2018-11-28 10:29:55','2018-11-28 10:29:55');`,
		"photos": `
INSERT INTO photos (id, url, moment_id, created_at, updated_at)
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
`,
		"moments": `
INSERT INTO moments (id, user_id, content, comment_total, like_total, status, created_at, updated_at)
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
`,
	}
)

func RunWithSchema(t *testing.T, test func(t *testing.T)) {
	db := Sqlx()
	db2 := Sqlx("db2")
	defer func() {
		for k := range createSchemas {
			_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", k))
			if err != nil {
				t.Error(err)
			}
		}
	}()

	for k, v := range createSchemas {
		udb := db
		if k == "photos" {
			udb = db2
		}
		_, err := udb.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", k))
		if err != nil {
			t.Error(err)
		}

		_, err = udb.Exec(v)
		if err != nil {
			t.Fatalf("create schema %s error:%s", k, err)
		}
	}

	test(t)
}

func initDatas(t *testing.T) {
	db := Sqlx()
	db2 := Sqlx("db2")
	for k, v := range datas {
		udb := db
		if k == "photos" {
			udb = db2
		}
		_, err := udb.Exec(v)
		if err != nil {
			t.Fatalf("init %s data error:%s", k, err)
		}
	}
}

func insert(id int) {
	user := &models.Users{
		Id:     id,
		Name:   "test" + strconv.Itoa(id),
		Status: 1,
	}
	_, err := Model(user).Create()
	if err != nil {
		log.Fatal(err)
	}
}

func insertStatus(id int, status int) {
	user := &models.Users{
		Id:     id,
		Name:   "test" + strconv.Itoa(id),
		Status: status,
	}
	Model(user).Create()
}

func TestBuilder_Get(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		{
			user := &models.Users{}
			err := Model(user).Where("id = ?", 1).Get()

			if err != nil {
				t.Error(err)
			}
			//fmt.Println(user)

		}

		{
			user := &models.Users{
				Name:   "test1",
				Status: 1,
			}
			err := Model(user).Get()

			if err != nil {
				t.Error(err)
			}
			fmt.Println(user)
		}

		{
			insertStatus(2, 0)
			user := &models.Users{
				Status: 0,
			}

			err := Model(user).Where("id = ?", 2).Get("status")

			if err != nil {
				t.Error(err)
			}
			fmt.Println(user)
		}
	})
}

func TestBuilder_Hint(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)

		user := make([]*models.Users, 0)
		err := Model(&user).Hint("/*+TDDL:slave()*/").All()

		if err != nil {
			t.Error(err)
		}

		fmt.Println(jsonEncode(user))
	})
}

func jsonEncode(i interface{}) string {
	ret, _ := json.Marshal(i)
	return string(ret)
}

func TestBuilder_All(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)

		user := make([]*models.Users, 0)
		err := Model(&user).All()

		if err != nil {
			t.Error(err)
		}

		fmt.Println(jsonEncode(user))
	})
}

func TestBuilder_Select(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)

		user := make([]*models.Users, 0)
		err := Model(&user).Select("id,name").All()

		if err != nil {
			t.Error(err)
		}

		fmt.Println(jsonEncode(user))
	})
}

func TestBuilder_InAll(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(3)
		insert(4)
		insert(5)

		user := make([]*models.Users, 0)
		err := Model(&user).Where("status = ? and id in(?)", 1, []int{1, 3, 4}).All()

		if err != nil {
			t.Error(err)
		}

		fmt.Println(jsonEncode(user))
	})
}

func TestBuilder_Update(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)

		{
			user := &models.Users{
				Name: "test2",
			}

			affected, err := Model(user).Where("id=?", 1).Update()

			if err != nil {
				t.Error("update user error", err)
			}

			if affected == 0 {
				t.Error("update user affected error", err)
			}
		}

		{
			user := &models.Users{
				Id:   1,
				Name: "test3",
			}

			affected, err := Model(user).Update()

			if err != nil {
				t.Error("update user error", err)
			}

			if affected == 0 {
				t.Error("update user affected error", err)
			}
		}
	})
}

func TestBuilder_Delete(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		{
			insert(1)
			affected, err := Model(&models.Users{}).Where("id=?", 1).Delete()

			if err != nil {
				t.Error("delete user error", err)
			}

			if affected == 0 {
				t.Error("delete user affected error", err)
			}
		}
		{
			insert(1)
			affected, err := Model(&models.Users{Id: 1}).Delete()

			if err != nil {
				t.Error("delete user error", err)
			}

			if affected == 0 {
				t.Error("delete user affected error", err)
			}
		}

		{
			insertStatus(1, 0)
			insertStatus(2, 0)
			insertStatus(3, 0)

			affected, err := Model(&models.Users{Status: 0}).Delete("status")

			if err != nil {
				t.Error("delete user error", err)
			}

			if affected != 3 {
				t.Error("delete user affected error", err)
			}
		}
	})
}

func TestBuilder_Count(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		{
			num, err := Model(&models.Users{}).Count()

			if err != nil {
				t.Error(err)
			}

			if num != 1 {
				t.Error("count user error")
			}
		}

		{
			insertStatus(2, 0)
			insertStatus(3, 0)

			num, err := Model(&models.Users{Status: 0}).Count("status")

			if err != nil {
				t.Error(err)
			}

			if num != 2 {
				t.Error("count user error")
			}
		}
	})
}

func TestBuilder_Create(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		user := &models.Users{
			//Id:    1,
			Name: "test",
		}
		id, err := Model(user).Create()

		if err != nil {
			t.Error(err)
		}

		if id != 1 {
			t.Error("lastInsertId error", id)
		}

		if int(id) != user.Id {
			t.Error("fill primaryKey error", id)
		}
	})
}

func TestBuilder_Limit(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(3)
		user := &models.Users{}
		err := Model(user).Limit(1).Get()

		if err != nil {
			t.Error(err)
		}
	})
}

func TestBuilder_Offset(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(3)
		user := &models.Users{}
		err := Model(user).Limit(1).Offset(1).Get()

		if err != nil {
			t.Error(err)
		}
	})
}

func TestBuilder_OrderBy(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(3)
		user := &models.Users{}
		err := Model(user).OrderBy("id desc").Limit(1).Offset(1).Get()

		if err != nil {
			t.Error(err)
		}

		if user.Id != 2 {
			t.Error("order by error")
		}

		//fmt.Println(user)
	})
}

func TestBuilder_Where(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(3)
		user := make([]*models.Users, 0)
		err := Model(&user).Where("id in(?,?)", 2, 3).OrderBy("id desc").All()

		if err != nil {
			t.Error(err)
		}

		if len(user) != 2 {
			t.Error("where error")
		}

		//fmt.Println(user)
	})
}

func TestBuilder_NullString(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		ct, _ := time.Parse("2006-01-02 15:04:05", "2018-09-02 00:00:00")
		{
			user := &models.Users{
				Id:     1,
				Name:   "test",
				Status: 1,
				SuccessTime: sql.NullString{
					String: "2018-09-03 00:00:00",
					Valid:  true,
				},
				CreatedAt: ct,
			}
			_, err := Model(user).Create()
			if err != nil {
				log.Fatal(err)
			}
		}

		{
			user := &models.Users{}
			err := Model(user).Where("id=1").Get()

			if err != nil {
				t.Error(err)
			}

			fmt.Println(jsonEncode(user))
		}

		{
			user := &models.Users{
				Id: 1,
				SuccessTime: sql.NullString{
					String: "2018-09-03 00:00:00",
					Valid:  true,
				},
				CreatedAt: ct,
			}

			err := Model(user).Get()

			if err != nil {
				t.Error(err)
			}

			fmt.Println(jsonEncode(user))
		}
	})
}

func TestBuilder_Relation1(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		initDatas(t)
		moment := &MomentList{}
		err := Model(moment).Relation("User", func(b *ModelStruct) {
			b.Where("status = 1")
		}).Where("status = 1 and id = ?", 14).Get()

		b, _ := json.MarshalIndent(moment, "", "	")
		fmt.Println(string(b), err)

		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestBuilder_Relation2(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		var moments = make([]*MomentList, 0)
		err := Model(&moments).Relation("User", func(b *ModelStruct) {
			b.Where("status = 0")
		}).Where("status = 1").Limit(10).All()

		b, _ := json.MarshalIndent(moments, "", "	")
		fmt.Println(string(b), err)

		if err != nil {
			t.Fatal(err)
		}
	})
}
