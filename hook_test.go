package gosql

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ilibs/gosql/v2/internal/example/models"
)

type hookUser struct {
	models.Users
}

func (u *hookUser) BeforeCreate() error {
	fmt.Println("BeforCreate run")
	if u.Id == 1 {
		return errors.New("before error")
	}

	return nil
}

func (u *hookUser) AfterCreate() {
	fmt.Println("AfterCreate run")
}

func (u *hookUser) BeforeUpdate() {
	fmt.Println("BeforeUpdate run")
}

func (u *hookUser) AfterUpdate() error {
	fmt.Println("AfterUpdate run")
	user := &models.Users{
		Id: 999,
	}

	err := Model(user).Get()
	return err
}

func (u *hookUser) BeforeDelete() {
	fmt.Println("BeforeDelete run")
}

func (u *hookUser) AfterDelete() {
	fmt.Println("AfterDelete run")
}

func (u *hookUser) AfterFind() {
	u.Name = "AfterUserName"
	fmt.Println("AfterFind run")
}

func TestNewHook(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		{
			user := &hookUser{models.Users{
				Id:     1,
				Name:   "test",
				Status: 1,
			}}
			_, err := Model(user).Create()
			if err == nil {
				t.Error("before create must error")
			}
		}

		{
			insert(2)
			user := &hookUser{models.Users{
				Id: 2,
			},
			}
			_, err := Model(user).Update()
			if err == nil {
				t.Error("after update must error")
			}
		}

		{
			user := &hookUser{models.Users{
				Id:     3,
				Name:   "test",
				Status: 1,
			},
			}
			_, err := Model(user).Create()
			if err != nil {
				t.Fatal(err)
			}

			user.Name = "test2"
			Model(user).Update()
			user2 := &hookUser{}
			Model(user2).Where("id=3").Get()
			if user2.Name != "AfterUserName" {
				t.Error("AfterFind change username error")
			}

			Model(user).Delete()
		}
	})
}

func TestHook_Err(t *testing.T) {
	hook := NewHook(nil)
	hook.Err(errors.New("test"))
	if hook.HasError() != 1 {
		t.Error("hook err")
	}
}

func TestHook_HasError(t *testing.T) {
	hook := NewHook(nil)
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	if hook.HasError() != 5 {
		t.Error("has error err")
	}
}

func TestHook_Error(t *testing.T) {
	hook := NewHook(nil)
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	if strings.Count(hook.Error().Error(), "test") != 5 {
		t.Error("get error err")
	}
}
