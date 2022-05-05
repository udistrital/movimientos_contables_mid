package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/movimientos_contables_mid/helpers/crud/movimientos_contables"
	e "github.com/udistrital/utils_oas/errorctrl"
)

// MovimientosController operations for Movimientos
type MovimientosController struct {
	beego.Controller
}

// URLMapping ...
func (c *MovimientosController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
}

// GetAll ...
// @Title GetAll
// @Description get Movimientos
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	detailfields	query	string	false	"Detailfields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Movimientos
// @Failure 403
// @router / [get]
func (c *MovimientosController) GetAll() {
	defer e.ErrorControlController(c.Controller, "MovimientosController")

	var fields []string
	var detailfields []string
	var sortby []string
	var order []string
	var query string
	var limit int = 10
	var offset int

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// detailfields: col1,col2,entity.col3
	if v := c.GetString("detailfields"); v != "" {
		detailfields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt("offset"); err == nil {
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
		query = v
	}
	var movimientos interface{}

	err := movimientos_contables.GetMovimientos(query, fields, limit, offset, sortby, order, detailfields, &movimientos)
	if err != nil {
		logs.Error(err)
		c.Data["mesaage"] = "Error service GetAll: The request contains an incorrect parameter or no record exists"
		c.Abort("404")
	} else {
		var data interface{}
		if movimientos != nil {
			data = movimientos
		} else {
			data = []interface{}{}
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": data}
	}
	c.ServeJSON()
}
