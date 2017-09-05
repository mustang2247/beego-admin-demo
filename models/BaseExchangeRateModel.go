package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"log"
	"reflect"
	"strings"
	"time"
)

//BaseExchangeRate 系统维护的汇率表
//	特别注意汇率的有效期
type BaseExchangeRate struct {
	Id int64 `orm:"column(id);pk;auto"`
	//汇率有效期，左右闭合原则
	StartAt string `orm:"column(start_at);size(8)"`
	OverAt  string `orm:"column(over_at);size(8)"`
	//源货币代号
	SrcHuobdaih string `orm:"column(srchuobdaih);size(2);valid(required)"`
	//目标货币代号
	DstHuobdaih string `orm:"column(dsthuobdaih);size(2);valid(required)"`
	//汇率
	ExchangeRate float64 `orm:"column(exchangerate);type(float64-decimal);decimals(6);digits(23);null"`

	//============= 与标签有有关结束 =============
	//管控部分
	Status   uint64    `orm:"column(status);default(2)" form:"Status"`
	Creator  *User     `orm:"column(creator);null;rel(fk)"`
	CreateAt time.Time `orm:"column(create_at);type(datetime);auto_now_add"`
	Updater  *User     `orm:"column(updater);null;rel(fk)"`
	UpdateAt time.Time `orm:"column(update_at);type(datetime);null"`
	Message  string    `orm:"column(message);size(255);null"`
}

// //节点表
// type BaseExchangeRate struct {
// 	SQLModel `sql:",inline"`
//
// 	Bizdate time.Time `orm:"unique;column(bizdate);type(datetime);null"`
// 	//@todo TODO DONE
// 	State string `orm:"unique;column(bizdate);type(datetime);null"	form:"State"  valid:"Required"`
// }

func (n *BaseExchangeRate) TableName() string {
	return beego.AppConfig.String("base_common_exchangerate_table")
}

//验证用户信息
func checkBaseExchangeRate(u *BaseExchangeRate) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&u)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

func init() {
	orm.RegisterModel(new(BaseExchangeRate))
}

//get node list
func GetBaseExchangeRatelist(page int64, page_size int64, sort string) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(BaseExchangeRate)
	qs := o.QueryTable(node)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}

	qs.Limit(page_size, offset).OrderBy(sort).Values(&nodes, "Id",
		"StartAt",
		"OverAt",
		"SrcHuobdaih",
		"DstHuobdaih",
		"ExchangeRate",
		"Status")
	count, _ = qs.Count()
	return nodes, count
}

func GetBaseExchangeRateById(nid int64) (BaseExchangeRate, error) {
	o := orm.NewOrm()
	node := BaseExchangeRate{Id: nid}
	err := o.Read(&node)
	if err != nil {
		return node, err
	}
	return node, nil
}

//添加用户
func AddBaseExchangeRate(n *BaseExchangeRate) (int64, error) {
	if err := checkBaseExchangeRate(n); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	node := new(BaseExchangeRate)
	// node.Name = n.Name
	// node.State = n.State
	// node.Country = n.Country

	id, err := o.Insert(node)
	return id, err
}

//更新用户
func UpdateBaseExchangeRate(n *BaseExchangeRate) (int64, error) {
	if err := checkBaseExchangeRate(n); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	node := make(orm.Params)
	//@todo 校验条件
	// if len(node) == 0 {
	// 	return 0, errors.New("update field is empty")
	// }
	var table BaseExchangeRate
	num, err := o.QueryTable(table).Filter("Id", n.Id).Update(node)
	return num, err
}

func DelBaseExchangeRateById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&BaseExchangeRate{Id: Id})
	return status, err
}

func GetBaseExchangeRatelistByGroupid(Groupid int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(BaseExchangeRate)
	count, _ = o.QueryTable(node).Filter("Group", Groupid).RelatedSel().Values(&nodes)
	return nodes, count
}

func GetBaseExchangeRateTree(pid int64, level int64) ([]orm.Params, error) {
	o := orm.NewOrm()
	node := new(BaseExchangeRate)
	var nodes []orm.Params
	_, err := o.QueryTable(node).Filter("Pid", pid).Filter("Level", level).Filter("Status", 2).Values(&nodes)
	if err != nil {
		return nodes, err
	}
	return nodes, nil
}

// GetAllDownloadtrac retrieves all Downloadtrac matches certain condition. Returns empty list if
// no records exist
func GetAllBaseExchangeRates(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Downloadtrac))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Downloadtrac
	qs = qs.OrderBy(sortFields...)
	if _, err := qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}
