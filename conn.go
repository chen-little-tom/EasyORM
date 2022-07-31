/**
 * @Author: chenchen
 * @Description:
 * @File: conn
 * @Version: 1.0.0
 * @Data: 2022/07/29 19:53
 */

package conn

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

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

// NewMysql 新建一个Mysql数据库连接
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

// resetEasyormEngine 重置引擎
func (e *EasyormEngine) resetEasyormEngine() *EasyormEngine {

	return e
}

// setErrorInfo 自定义错误格式
func (e *EasyormEngine) setErrorInfo(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return errors.New("File: " + file + ":" + strconv.Itoa(line) + "," + err.Error())
}

func (e *EasyormEngine) insertData(data interface{}, insertType string) (int64, error) {

	//反射type和value
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	//字段名
	var fieldName []string

	//问号?占位符
	var placeholder []string

	//循环判断
	for i := 0; i < t.NumField(); i++ {

		//小写开头，无法反射，跳过
		if !v.Field(i).CanInterface() {
			continue
		}

		//解析tag，找出真实的sql字段名
		sqlTag := t.Field(i).Tag.Get("sql")
		if sqlTag != "" {
			//跳过自增字段
			if strings.Contains(strings.ToLower(sqlTag), "auto_increment") {
				continue
			} else {
				fieldName = append(fieldName, strings.Split(sqlTag, ",")[0])
				placeholder = append(placeholder, "?")
			}
		} else {
			fieldName = append(fieldName, t.Field(i).Name)
			placeholder = append(placeholder, "?")
		}

		//字段的值
		e.AllExec = append(e.AllExec, v.Field(i).Interface())
	}

	//拼接表，字段名，占位符
	e.Prepare = insertType + " into " + e.GetTable() + " (" + strings.Join(fieldName, ",") + ") values(" + strings.Join(placeholder, ",") + ")"

	//prepare
	var stmt *sql.Stmt
	var err error
	stmt, err = e.Db.Prepare(e.Prepare)
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//执行exec，注意这里是stmt.Exec
	result, err := stmt.Exec(e.AllExec...)
	//方法运行完之后清空AllExec中保存的数据
	defer func() {
		e.AllExec = nil
	}()
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//获取自增ID
	id, _ := result.LastInsertId()
	return id, nil
}

// batchInsertData 批量插入数据
func (e *EasyormEngine) batchInsertData(batchData interface{}, insertType string) (int64, error) {

	//反射解析
	getValue := reflect.ValueOf(batchData)

	//切片大小
	l := getValue.Len()

	//字段名
	var fieldName []string

	//占位符
	var placeholderString []string

	//循环判断
	for i := 0; i < l; i++ {
		value := getValue.Index(i)
		typed := value.Type()
		if typed.Kind() != reflect.Struct {
			panic("Subelements inserted in batches must be of structure type")
		}

		num := value.NumField()

		//子元素的值
		var placeholder []string
		//循环遍历子元素
		for j := 0; j < num; j++ {

			//小写开头，无法反射，跳过
			if !value.Field(j).CanInterface() {
				continue
			}

			//解析tag，找出真实的sql字段名
			sqlTag := typed.Field(j).Tag.Get("sql")
			if sqlTag != "" {
				//跳过自增字段
				if strings.Contains(strings.ToLower(sqlTag), "auto_increment") {
					continue
				} else {
					//字段名只记录第一个的
					if i == 1 {
						fieldName = append(fieldName, strings.Split(sqlTag, ",")[0])
					}
					placeholder = append(placeholder, "?")
				}
			} else {
				//字段名只记录第一个的
				if i == 1 {
					fieldName = append(fieldName, typed.Field(j).Name)
				}
				placeholder = append(placeholder, "?")
			}
			//字段值
			e.AllExec = append(e.AllExec, value.Field(j).Interface())
		}

		//子元素拼接成多个()括号后的值
		placeholderString = append(placeholderString, "("+strings.Join(placeholder, ",")+")")
	}

	//拼接表、字段名、占位符
	e.Prepare = insertType + " into " + e.GetTable() + " (" + strings.Join(fieldName, ",") + ") values " +
		strings.Join(placeholderString, ",")

	//prepare
	var stmt *sql.Stmt
	var err error
	stmt, err = e.Db.Prepare(e.Prepare)
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//执行exec，注意这里是stmt.Exec
	result, err := stmt.Exec(e.AllExec...)
	//方法运行完之后清空AllExec中保存的数据
	defer func() {
		e.AllExec = nil
	}()
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//获取自增ID
	id, _ := result.LastInsertId()
	return id, nil
}
