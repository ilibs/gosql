package gosql

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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

	err := WithContext(nil).Model(user).Get()
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
			_, err := WithContext(nil).Model(user).Create()
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
	hook := NewHook(nil, nil)
	hook.Err(errors.New("test"))
	if !hook.HasError() {
		t.Error("hook err")
	}
}

func TestHook_HasError(t *testing.T) {
	hook := NewHook(nil, nil)
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	if !hook.HasError() {
		t.Error("has error err")
	}
}

func TestHook_Error(t *testing.T) {
	hook := NewHook(nil, nil)
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	hook.Err(errors.New("test"))
	if strings.Count(hook.Error().Error(), "test") != 5 {
		t.Error("get error err")
	}
}

type testModelCallBack struct {
}

func (m *testModelCallBack) BeforeCreate() {
}

func (m *testModelCallBack) AfterCreate() error {
	return nil
}

func (m *testModelCallBack) BeforeChange(tx *DB) {
}

func (m *testModelCallBack) AfterChange(tx *DB) error {
	return nil
}

func (m *testModelCallBack) BeforeUpdate(ctx context.Context) {
}

func (m *testModelCallBack) AfterUpdate(ctx context.Context) error {
	return nil
}

func (m *testModelCallBack) BeforeDelete(ctx context.Context, tx *DB) {
}

func (m *testModelCallBack) AfterDelete(ctx context.Context, tx *DB) error {
	return nil
}

func TestHook_callMethod(t *testing.T) {

	hook := NewHook(nil, nil)

	m := &testModelCallBack{}

	refVal := reflect.ValueOf(m)
	hook.callMethod("BeforeCreate", refVal)
	hook.callMethod("BeforeChange", refVal)
	hook.callMethod("BeforeDelete", refVal)
	hook.callMethod("BeforeUpdate", refVal)
	hook.callMethod("AfterCreate", refVal)
	hook.callMethod("AfterChange", refVal)
	hook.callMethod("AfterDelete", refVal)
	hook.callMethod("AfterUpdate", refVal)
}
