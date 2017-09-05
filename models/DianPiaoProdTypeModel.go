package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"log"
	"reflect"
	"strings"
	"time"
)

// DianPiaoProductType 电票产品类型
type DianPiaoProductType struct {
	Id              int64  `orm:"column(id);auto;pk"`
	ProductTypeCode string `orm:"column(producttype_code);unique;size(2)" form:"ProductTypeCode" valid:"Required"`
	Description     string `orm:"column(description);size(255)" form:"Description"  valid:"Required"`

	//管控部分
	Status   uint64    `orm:"column(status);default(2)" form:"Status"`
	Creator  *User     `orm:"column(creator);null;rel(fk)"`
	CreateAt time.Time `orm:"column(create_at);type(datetime);auto_now_add"`
	Updater  *User     `orm:"column(updater);null;rel(fk)"`
	UpdateAt time.Time `orm:"column(update_at);type(datetime);null"`
	Message  string    `orm:"column(message);size(255);null"`
}

func (n *DianPiaoProductType) TableName() string {
	return beego.AppConfig.String("base_dianpiaoproducttype_table")
}

//验证用户信息
func checkDianPiaoProductType(u *DianPiaoProductType) (err error) {
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
	orm.RegisterModel(new(DianPiaoProductType))
}

//get node list
func GetDianPiaoProductTypelist(page int64, page_size int64, sort string) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(DianPiaoProductType)
	qs := o.QueryTable(node)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).Values(&nodes, "Id", "DianPiaoProductType", "Description", "Cabinet", "Folder", "FileType", "Paramdesc", "ETLTaskDescId", "UpdateAt", "CreateAt")
	count, _ = qs.Count()
	return nodes, count
}

func GetDianPiaoProductTypeById(nid int64) (DianPiaoProductType, error) {
	o := orm.NewOrm()
	node := DianPiaoProductType{Id: nid}
	err := o.Read(&node)
	if err != nil {
		return node, err
	}
	return node, nil
}

//添加用户
func AddDianPiaoProductType(n *DianPiaoProductType) (int64, error) {
	if err := checkDianPiaoProductType(n); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	node := new(DianPiaoProductType)
	node.ProductTypeCode = n.ProductTypeCode
	node.Description = n.Description

	id, err := o.Insert(node)
	return id, err
}

//更新用户
func UpdateDianPiaoProductType(n *DianPiaoProductType) (int64, error) {
	if err := checkDianPiaoProductType(n); err != nil {
		return 0, err
	}
	fmt.Printf("UpdateDianPiaoProductType.DianPiaoProductType %q\n", n)
	o := orm.NewOrm()
	node := make(orm.Params)
	if len(n.ProductTypeCode) > 0 {
		node["ProductTypeCode"] = n.ProductTypeCode
	}
	// if len(n.Name) > 0 {
	// 	node["Name"] = n.Name
	// }
	// if len(n.Remark) > 0 {
	// 	node["Remark"] = n.Remark
	// }
	// if n.Status != 0 {
	// 	node["Status"] = n.Status
	// }
	if len(node) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table DianPiaoProductType
	num, err := o.QueryTable(table).Filter("Id", n.Id).Update(node)
	return num, err
}

func DelDianPiaoProductTypeById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&DianPiaoProductType{Id: Id})
	return status, err
}

func GetDianPiaoProductTypelistByGroupid(Groupid int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(DianPiaoProductType)
	count, _ = o.QueryTable(node).Filter("Group", Groupid).RelatedSel().Values(&nodes)
	return nodes, count
}

func GetDianPiaoProductTypeTree(pid int64, level int64) ([]orm.Params, error) {
	o := orm.NewOrm()
	node := new(DianPiaoProductType)
	var nodes []orm.Params
	_, err := o.QueryTable(node).Filter("Pid", pid).Filter("Level", level).Filter("Status", 2).Values(&nodes)
	if err != nil {
		return nodes, err
	}
	return nodes, nil
}

// GetAllDianPiaoProductType retrieves all DianPiaoProductType matches certain condition. Returns empty list if
// no records exist
func GetAllDianPiaoProductTypes(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(DianPiaoProductType))
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

	var l []DianPiaoProductType
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
