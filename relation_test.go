package gosql

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ilibs/gosql/example/models"
)

type UserMoment struct {
	models.Moments
	User   *models.Users    `json:"user" db:"-" relation:"user_id,id"`
	Photos []*models.Photos `json:"photos" db:"-" relation:"id,moment_id"`
}

func TestRelationOne1(t *testing.T) {
	moment := &UserMoment{}
	err := Model(moment).Where("status = 1 and id = ?",14).Get()

	b , _ :=json.MarshalIndent(moment,"","	")
	fmt.Println(string(b), err)

	if err != nil {
		t.Fatal(err)
	}

	if moment.User.NickName == "" {
		t.Fatal("relation one-to-one data error[user]")
	}

	if len(moment.Photos) == 0 {
		t.Fatal("relation get one-to-many data error[photos]")
	}
}


func TestRelationOne2(t *testing.T) {
	moment := &UserMoment{}
	err := Model(moment).Relation("User" , func(b *Builder) *Builder {
		return b.Where("id = 2")
	}).Where("status = 1 and id = ?",14).Get()

	b , _ :=json.MarshalIndent(moment,"","	")
	fmt.Println(string(b), err)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRelationOne3(t *testing.T) {
	moment := &UserMoment{}
	err := Relation("User" , func(b *Builder) *Builder {
		return b.Where("id = 2")
	}).Get(moment , "select * from moments")

	b , _ :=json.MarshalIndent(moment,"","	")
	fmt.Println(string(b), err)

	if err != nil {
		t.Fatal(err)
	}
}



func TestRelationAll1(t *testing.T) {
	var moments = make([]*UserMoment, 0)
	err := Model(&moments).Where("status = 1").Limit(10).All()
	if err != nil {
		t.Fatal(err)
	}


	b , _ :=json.MarshalIndent(moments,"","	")
	fmt.Println(string(b),err)

	if len(moments) == 0 {
		t.Fatal("relation get many-to-many data error[moments]")
	}

	if moments[0].User.NickName == "" {
		t.Fatal("relation get many-to-many data error[user]")
	}

	if len(moments[0].Photos) == 0 {
		t.Fatal("relation get many-to-many data error[photos]")
	}

}

func TestRelationAll2(t *testing.T) {
	var moments = make([]*UserMoment, 0)
	err := Model(&moments).Relation("User"  , func(b *Builder) *Builder {
		return b.Where("id = 2")
	}).Where("status = 1").Limit(10).All()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRelationAll3(t *testing.T) {
	var moments = make([]*UserMoment, 0)
	err := Relation("User"  , func(b *Builder) *Builder {
		return b.Where("id = 2")
	}).Select(&moments , "select * from moments")

	if err != nil {
		t.Fatal(err)
	}
}