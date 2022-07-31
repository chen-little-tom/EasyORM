/**
 * @Author: chenchen
 * @Description:
 * @File: ormFunc_test
 * @Version: 1.0.0
 * @Data: 2022/07/30 18:22
 */

package conn

import (
	"fmt"
	"testing"
)

type User struct {
	Name string `sql:"name"`
	Age  int    `sql:"age"`
}

type UserWithID struct {
	Id   int64  `sql:"id, auto_increment"`
	Name string `sql:"name"`
	Age  int    `sql:"age"`
}

func TestEasyormEngine_Insert(t *testing.T) {
	conn, _ := NewMysql("root", "chenchen521", "localhost", "orm")
	user1 := User{
		Name: "cc",
		Age:  25,
	}
	conn.Table("TABLE_USER").Insert(user1)

	user2 := UserWithID{
		Name: "cc",
		Age:  24,
	}
	id, _ := conn.Table("TABLE_USER_ID").Insert(user2)
	user2.Id = id
	fmt.Printf("Auto_increment Id is %d", user2.Id)
}
