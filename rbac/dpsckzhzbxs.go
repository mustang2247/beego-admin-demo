package rbac

import (
	"encoding/json"
	"errors"
	"github.com/devuser/beego-admin-demo/models"
	"strconv"
	"strings"

	// "github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// oprations for DPSCkzhzbx
type DPSCkzhzbxController struct {
	CommonController
}

func (c *DPSCkzhzbxController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// @Title Post
// @Description create DPSCkzhzbx
// @Param	body		body 	models.DPSCkzhzbx	true		"body for DPSCkzhzbx content"
// @Success 201 {int64} models.DPSCkzhzbx
// @Failure 403 body is empty
// @router / [post]
func (c *DPSCkzhzbxController) Post() {
	var v models.DPSCkzhzbx
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddDPSCkzhzbx(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// @Title Get
// @Description get DPSCkzhzbx by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DPSCkzhzbx
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DPSCkzhzbxController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	DPSCkzhzbxId, _ := strconv.Atoi(idStr)
	v, err := models.GetDPSCkzhzbxById(int64(DPSCkzhzbxId))
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// @Title Get All
// @Description get DPSCkzhzbx
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an int64eger"
// @Param	offset	query	string	false	"Start position of result set. Must be an int64eger"
// @Success 200 {object} models.DPSCkzhzbx
// @Failure 403
// @router / [get]
func (c *DPSCkzhzbxController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query map[string]string = make(map[string]string)
	var limit int64 = 10
	var offset int64 = 0

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.Split(cond, ":")
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllDPSCkzhzbxs(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// @Title Update
// @Description update the DPSCkzhzbx
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.DPSCkzhzbx	true		"body for DPSCkzhzbx content"
// @Success 200 {object} models.DPSCkzhzbx
// @Failure 403 :id is not int64
// @router /:id [put]
func (c *DPSCkzhzbxController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DPSCkzhzbx{Id: int64(id)}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateDPSCkzhzbxById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// @Title Delete
// @Description delete the DPSCkzhzbx
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *DPSCkzhzbxController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	DPSCkzhzbxId, _ := strconv.Atoi(idStr)
	if err := models.DelDPSCkzhzbxById(int64(DPSCkzhzbxId)); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (this *DPSCkzhzbxController) Index() {
	if this.IsAjax() {
		page, _ := this.GetInt64("page")
		page_size, _ := this.GetInt64("rows")
		sort := this.GetString("sort")
		order := this.GetString("order")
		if len(order) > 0 {
			if order == "desc" {
				sort = "-" + sort
			}
		} else {
			sort = "Id"
		}
		roles, count := models.GetDPSCkzhzbxList(page, page_size, sort)
		//GetAllDPSCkzhzbxs(page, page_size, sort)
		if len(roles) < 1 {
			roles = []orm.Params{}
		}
		this.Data["json"] = &map[string]interface{}{"total": count, "rows": &roles}
		this.ServeJSON()
		return
	} else {
		this.TplName = this.GetTemplatetype() + "/rbac/dpsdetail.tpl"
	}

}
