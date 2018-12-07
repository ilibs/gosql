package gosql

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ilibs/gosql/example/models"
)


type UserMoment struct {
	models.Users
	Moments   []*models.Moments    `json:"moments" db:"-" relation:"id,user_id"`
}

func TestRelationOne2(t *testing.T) {
	moment := &UserMoment{}
	err := Model(moment).Relation("Moments", func(b *Builder) {
		b.Limit(2)
	}).Where("id = ?",5).Get()

	b , _ :=json.MarshalIndent(moment,"","	")
	fmt.Println(string(b), err)

	if err != nil {
		t.Fatal(err)
	}
}

type MomentList struct {
	models.Moments
	User   *models.Users    `json:"user" db:"-" relation:"user_id,id"`
	Photos []*models.Photos `json:"photos" db:"-" relation:"id,moment_id"`
}

func TestRelationOne(t *testing.T) {
	moment := &MomentList{}
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

func TestRelationAll(t *testing.T) {
	var moments = make([]*MomentList, 0)
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