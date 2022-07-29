/**
 * @Author: chenchen
 * @Description:
 * @File: conn
 * @Version: 1.0.0
 * @Data: 2022/07/29 19:53
 */

package conn

import "database/sql"

type EasyormEngine struct {
	Db           *sql.DB
	TableName    string
	Prepare      string
	AllExec      []interface{}
	Sql          string
	WhereParam   string
	LimitParam   string
	OrderParam   string
	OrWhereParam string
	WhereExec    []interface{}
	UpdateParam  string
	UpdateExec   []interface{}
	FieldParam   string
	TransStatus  int
	Tx           *sql.Tx
	GroupParam   string
	HavingParam  string
}

// resetEasyormEngine 重置引擎
func (e *EasyormEngine) resetEasyormEngine() *EasyormEngine {

	return e
}

func NewMysql(Username string, Password string, Address string, Dbname string) (*EasyormEngine, error) {
	dsn := Username + ":" + Password + "@tcp(" + Address + ")/" + Dbname + "?charset=utf8&timeout=5s&readTimeout=6s"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &EasyormEngine{
		Db:         db,
		FieldParam: "*",
	}, nil
}
