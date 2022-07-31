/**
 * @Author: chenchen
 * @Description:
 * @File: ormFunc
 * @Version: 1.0.0
 * @Data: 2022/07/29 20:08
 */

package conn

import (
	"errors"
	"reflect"
)

// Table 设置表名
func (e *EasyormEngine) Table(name string) *EasyormEngine {
	e.TableName = name

	//重置引擎
	e.resetEasyormEngine()
	return e
}

// GetTable 获取表名
func (e *EasyormEngine) GetTable() string {
	return e.TableName
}

// BatchInsert 批量插入
func (e *EasyormEngine) BatchInsert(data interface{}) (int64, error) {
	return e.batchInsertData(data, "insert")
}

// BatchReplace 批量替换
func (e *EasyormEngine) BatchReplace(data interface{}) (int64, error) {
	return e.batchInsertData(data, "replace")
}

func (e *EasyormEngine) Insert(data interface{}) (int64, error) {

	//判断是批量还是单个插入
	getValue := reflect.ValueOf(data).Kind()
	if getValue == reflect.Struct {
		return e.insertData(data, "insert")
	} else if getValue == reflect.Slice || getValue == reflect.Array {
		return e.batchInsertData(data, "insert")
	} else {
		return 0, errors.New("the inserted data format is incorrect. " +
			"the single table insert format is: struct, and the batch insert format is: []struct")
	}
}

func (e *EasyormEngine) Replace(data interface{}) (int64, error) {
	//判断是批量还是单个插入
	getValue := reflect.ValueOf(data).Kind()
	if getValue == reflect.Struct {
		return e.insertData(data, "replace")
	} else if getValue == reflect.Slice || getValue == reflect.Array {
		return e.batchInsertData(data, "replace")
	} else {
		return 0, errors.New("the replaced data format is incorrect. " +
			"the single table replace format is: struct, and the batch replace format is: []struct")
	}
}
