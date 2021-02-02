package gosql

import (
	"testing"

	"github.com/ilibs/gosql/v2/internal/example/models"
)

type UserMomentRelation struct {
	models.Users
	Moments []*models.Moments `json:"moments" db:"-" relation:"id,user_id" connection:"db2"`
}

func TestRelationOneWithRelationDB(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		initDatas(t)
		moment := &UserMomentRelation{}
		// Use("default").Model()
		err := Use("default").Model(NewModelWrapper(map[string]*DB{
			"default": Use("default"),
			"db2":     Use("db2"),
		}, moment)).Relation("Moments", func(b *Builder) {
			b.Limit(2)
		}).Where("id = ?", 5).Get()

		// b, _ := json.MarshalIndent(moment, "", "	")
		// fmt.Println(string(b), err)

		if err != nil {
			t.Fatal(err)
		}
	})
}

type UserMoment struct {
	models.Users
	Moments []*models.Moments `json:"moments" db:"-" relation:"id,user_id" connection:"db2"`
}

func TestRelationOne2(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		initDatas(t)
		moment := &UserMoment{}
		err := Model(moment).Relation("Moments", func(b *Builder) {
			b.Limit(2)
		}).Where("id = ?", 5).Get()

		// b, _ := json.MarshalIndent(moment, "", "	")
		// fmt.Println(string(b), err)

		if err != nil {
			t.Fatal(err)
		}
	})
}

type MomentList struct {
	models.Moments
	User   *models.Users    `json:"user" db:"-" relation:"user_id,id"`
	Photos []*models.Photos `json:"photos" db:"-" relation:"id,moment_id" connection:"db2"`
}

func TestRelationOne(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		initDatas(t)

		moment := &MomentList{}
		err := Model(moment).Where("status = 1 and id = ?", 14).Get()

		// b, _ := json.MarshalIndent(moment, "", "	")
		// fmt.Println(string(b), err)

		if err != nil {
			t.Fatal(err)
		}

		if moment.User.Name == "" {
			t.Fatal("relation one-to-one data error[user]")
		}

		if len(moment.Photos) == 0 {
			t.Fatal("relation get one-to-many data error[photos]")
		}
	})
}

func TestRelationAll(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		initDatas(t)

		var moments = make([]*MomentList, 0)
		err := Model(&moments).Where("status = 1").Limit(10).All()
		if err != nil {
			t.Fatal(err)
		}

		// b, _ := json.MarshalIndent(moments, "", "	")
		// fmt.Println(string(b), err)

		if len(moments) == 0 {
			t.Fatal("relation get many-to-many data error[moments]")
		}

		if moments[0].User.Name == "" {
			t.Fatal("relation get many-to-many data error[user]")
		}

		if len(moments[0].Photos) == 0 {
			t.Fatal("relation get many-to-many data error[photos]")
		}
	})
}

type MomentListWrapper struct {
	models.Moments
	User   *models.Users    `json:"user" db:"-" relation:"user_id,id"`
	Photos []*models.Photos `json:"photos" db:"-" relation:"id,moment_id" connection:"db2"`
}

func TestRelationModelWrapper(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		initDatas(t)
		var moments = make([]*MomentListWrapper, 0)
		err := Use("default").Model(NewModelWrapper(map[string]*DB{
			"default": Use("default"),
			"db2":     Use("db2"),
		}, &moments)).Where("status = 1").Limit(10).All()
		if err != nil {
			t.Fatal(err)
		}

		// b, _ := json.MarshalIndent(moments, "", "	")
		// fmt.Println(string(b), err)

		if len(moments) == 0 {
			t.Fatal("relation get many-to-many data error[moments]")
		}

		if moments[0].User.Name == "" {
			t.Fatal("relation get many-to-many data error[user]")
		}

		if len(moments[0].Photos) == 0 {
			t.Fatal("relation get many-to-many data error[photos]")
		}
	})
}
