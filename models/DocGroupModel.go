package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"reflect"
	"strings"
	"time"
)

type DocGroup struct {
	Id          int64  `orm:"column(id);auto;pk"`
	Groupname   string `orm:"column(groupname);unique;size(255)"`
	Description string `orm:"column(description);size(255)"`

	UserList []*User `orm:"rel(m2m)"`

	//管控部分
	Status   uint64    `orm:"column(status);default(2)" form:"Status"`
	Creator  *User     `orm:"column(creator);null;rel(fk)"`
	CreateAt time.Time `orm:"column(create_at);type(datetime);auto_now_add"`
	Updater  *User     `orm:"column(updater);null;rel(fk)"`
	UpdateAt time.Time `orm:"column(update_at);type(datetime);null"`
	Message  string    `orm:"column(message);size(255);null"`
}

func (t *DocGroup) TableName() string {
	return beego.AppConfig.String("doc_docgroup_table")
}

func init() {
	orm.RegisterModel(new(DocGroup))
}

// AddDocGroup insert a new DocGroup into database and returns
// last inserted Id on success.
func AddDocGroup(m *DocGroup) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetDocGroupById retrieves DocGroup by Id. Returns error if
// Id doesn't exist
func GetDocGroupById(id int64) (v *DocGroup, err error) {
	o := orm.NewOrm()
	v = &DocGroup{Id: (id)}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//get node list
func GetDocgrouplist(page int64, page_size int64, sort string) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(DocGroup)
	qs := o.QueryTable(node)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).Values(&nodes, "Id", "Groupname", "Description", "CreateAt", "UpdateAt")
	count, _ = qs.Count()
	return nodes, count
}

// GetAllDocGroup retrieves all DocGroup matches certain condition. Returns empty list if
// no records exist
func GetAllDocGroup(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(DocGroup))
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

	var l []DocGroup
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

// UpdateDocGroup updates DocGroup by Id and returns error if
// the record to be updated doesn't exist
func UpdateDocGroupById(m *DocGroup) (err error) {
	o := orm.NewOrm()
	v := DocGroup{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDocGroup deletes DocGroup by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDocGroup(id int64) (err error) {
	o := orm.NewOrm()
	v := DocGroup{Id: (id)}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&DocGroup{Id: (id)}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
