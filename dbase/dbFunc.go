package dbase

import (
	"context"
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

var ErrorFieldsIllegal = errors.New("[dbase]Query fields must be 0 or a multiple of 2")

//校验filter是否符合规定
func verifyFields(fields []interface{}) error {
	if len(fields)%2 != 0 {
		return ErrorFieldsIllegal
	}
	return nil
}

// List 用于返回db查询的一组数据，以传入数组ptr方式获取查询返回值
// in条件field key自行添加 __in, val为数组
func (d *DB) List(m interface{}, i interface{}, orders *[]string, cols *[]string, fields ...interface{}) (rows int64, err error) {
	err = verifyFields(fields)
	if err != nil {
		return 0, err
	}
	qs := d.QueryTable(m)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	if orders != nil {
		qs = qs.OrderBy(*orders...)
	}
	qs = qs.Limit(-1)
	if cols != nil {
		qs.All(i, *cols...)
	}
	return qs.All(i)
}

// ListRaw 用于返回db查询的一组数据，以[]orm.Params形式返回
// in条件field key自行添加 __in, val为数组
func (d *DB) ListRaw(m interface{}, orders *[]string, fields ...interface{}) (rows int64, data []orm.Params, err error) {
	err = verifyFields(fields)
	if err != nil {
		return 0, nil, err
	}
	qs := d.QueryTable(m)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	if orders != nil {
		qs = qs.OrderBy(*orders...)
	}
	qs = qs.Limit(-1)
	rows, err = qs.Values(&data)
	return
}

// One 用于返回db查询List的第一条数据，以传入数组ptr方式获取查询返回值
func (d *DB) One(m interface{}, i interface{}, orders *[]string, fields ...interface{}) (err error) {
	err = verifyFields(fields)
	if err != nil {
		return err
	}
	qs := d.QueryTable(m)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	if orders != nil {
		qs = qs.OrderBy(*orders...)
	}
	qs = qs.Limit(-1)
	err = qs.One(i)
	return
}

// Get 用于查询一条数据，以传入数组ptr方式获取查询返回值
func (d *DB) Get(m interface{}, cols ...string) (err error) {
	return d.ReadWithCtx(context.Background(), m, cols...)
}

// Insert 插入一条数据
func (d *DB) Insert(i interface{}) (id int64, err error) {
	return d.InsertWithCtx(context.Background(), i)
}

// InsertMulti 一次插入多条条数据
// perIns 单次插入数量
func (d *DB) InsertMulti(i interface{}, perIns int) (id int64, err error) {
	return d.InsertMultiWithCtx(context.Background(), perIns, i)
}

// Update 按字段更新数据
func (d *DB) Update(i interface{}, clos ...string) (rows int64, err error) {
	return d.UpdateWithCtx(context.Background(), i, clos...)
}

// UpgradeFilter 按字段更新数据
// filter 条件
/* values Set = v 条件
orm.Params{
    "name": "astaxie",
}
或
orm.Params{
    "nums": orm.ColValue(orm.ColAdd, 100),
}
ColAdd      // 加
ColMinus    // 减
ColMultiply // 乘
ColExcept   // 除
*/
func (d *DB) UpgradeFilter(i interface{}, filters *map[string]interface{}, values *orm.Params) (rows int64, err error) {
	qs := d.QueryTable(i)
	if filters != nil {
		for k, v := range *filters {
			qs = qs.Filter(k, v)
		}
	}
	return qs.Update(*values)
}

// InsertOrUpdate 插入或更新一条数据
func (d *DB) InsertOrUpdate(i interface{}, fields ...string) (rows int64, err error) {
	return d.InsertOrUpdateWithCtx(context.Background(), i, fields...)
}

// Delete 删除数据
func (d *DB) Delete(i interface{}, fields ...string) (rows int64, err error) {
	return d.DeleteWithCtx(context.Background(), i, fields...)
}

// Count 数量
func (d *DB) Count(i interface{}, fields ...interface{}) (count int64, err error) {
	err = verifyFields(fields)
	if err != nil {
		return
	}
	qs := d.QueryTable(i)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	count, err = qs.Count()
	return
}
