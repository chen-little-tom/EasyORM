/**
 * @Author: chenchen
 * @Description:
 * @File: ormFunc
 * @Version: 1.0.0
 * @Data: 2022/07/29 20:08
 */

package conn

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
