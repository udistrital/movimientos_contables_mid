package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/udistrital/movimientos_contables_mid/helpers/crud/cuentas_contables"
)

// CuentaContableController operations for Cuenta_contable
type CuentaContableController struct {
	beego.Controller
}

// URLMapping ...
func (c *CuentaContableController) URLMapping() {
	c.Mapping("Delete", c.Delete)
}

// Delete ...
// @Title Delete
// @Description delete the Cuenta_contable
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *CuentaContableController) Delete() {

	idStr := c.Ctx.Input.Param(":id")
	if err := cuentas_contables.DeleteNodoCuentaContableByCuentaId(idStr); err == nil {
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "OK"}
	} else {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"Fail": false, "Status": "500", "Message": "Fail", "Data": "fail"}
	}
	c.ServeJSON()

}
