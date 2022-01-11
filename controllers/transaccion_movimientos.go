package controllers

import (
	"encoding/json"
	_ "encoding/json"
	_ "fmt"
	_ "time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/astaxie/beego/logs"
	"github.com/udistrital/movimientos_contables_mid/helpers"
	_ "github.com/udistrital/movimientos_contables_mid/helpers"
	"github.com/udistrital/movimientos_contables_mid/models"
)

// TransaccionMovimientosController operations for TransaccionMovimientos
type TransaccionMovimientosController struct {
	beego.Controller
}

//URLMapping ...
func (c *TransaccionMovimientosController) URLMapping() {
	c.Mapping("PostTransaccionMovimientos", c.PostTransaccionMovimientos)
}

// PostTransaccionMovimientosController ...
// @Title PostTransaccionMovimientos
// @Description create PostTransaccionMovimientos - DEPRECADA: Se eliminará a futuro. Por favor usar el endpoint sencillo, sin el `transaccion_movimientos` adicional
// @Deprecated deprecated
// @Param	body		body 	models.TransaccionMovimientos	true		"body for TransaccionMovimientos content"
// @Success 201			Ok
// @Failure 400 the request contains incorrect syntax
// @Failure 500 Unhandled Error
// @router /transaccion_movimientos [post]
func (c *TransaccionMovimientosController) PostTransaccionMovimientosDeprecada() {
	c.PostTransaccionMovimientos()
}

// PostTransaccionMovimientosController ...
// @Title PostTransaccionMovimientos
// @Description create PostTransaccionMovimientos
// @Param	body		body 	models.TransaccionMovimientos	true		"body for TransaccionMovimientos content"
// @Success 201			Ok
// @Failure 400 the request contains incorrect syntax
// @Failure 500 Unhandled Error
// @router / [post]
func (c *TransaccionMovimientosController) PostTransaccionMovimientos() {
	//defer helpers.GestionError(c)
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "TransaccionMovimientosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500")
			}
		}
	}()

	var v models.TransaccionMovimientos = models.TransaccionMovimientos{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := helpers.RegistroTransaccionMovimientos(v); err == nil {
			//c.Data["json"] = "OK"
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": "OK"}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "PostTransaccionMovimientos", "err": err.Error(), "status": "400"})
	}

	c.ServeJSON()

}
