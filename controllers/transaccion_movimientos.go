package controllers

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "fmt"
	"net/http"
	"strconv"
	_ "time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/astaxie/beego/logs"
	e "github.com/udistrital/utils_oas/errorctrl"

	"github.com/udistrital/movimientos_contables_mid/helpers"
	_ "github.com/udistrital/movimientos_contables_mid/helpers"
	"github.com/udistrital/movimientos_contables_mid/helpers/transaccionmovimientos"
	"github.com/udistrital/movimientos_contables_mid/models"
)

// TransaccionMovimientosController operations for TransaccionMovimientos
type TransaccionMovimientosController struct {
	beego.Controller
}

//URLMapping ...
func (c *TransaccionMovimientosController) URLMapping() {
	c.Mapping("Post", c.PostTransaccionMovimientos)
}

// PostTransaccionMovimientosDeprecada ...
// @Title PostTransaccionMovimientosDeprecada
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

// Get ...
// @Title Get TransaccionMovimientos
// @Description Get TransaccionMovimientos
// @Param idType   path  string true  "buscar por id de: consecutivo o transaccion"
// @Param id       path  int    true  "El ID como tal"
// @Param detailed query bool   false "Traer los movimientos asociados? `false` por defecto"
// @Success 200 {object} models.TransaccionMovimientos
// @Failure 400 Parametros Incorrectos
// @Failure 404 Transaccion no encontrada
// @Failure 500 Error no manejado!
// @Failure 502 Error al contactar otra API
// @router /:idType/:id [get]
func (c *TransaccionMovimientosController) Get() {
	defer e.ErrorControlController(c.Controller, "TransaccionMovimientosController")
	const funcion string = "Get"

	var (
		detailed bool
		idType   string
		id       int
		err      error
	)
	if idType = c.GetString(":idType"); len(idType) == 0 {
		err := fmt.Errorf("error: Criterio no especificado")
		panic(e.Error(funcion, err, strconv.Itoa(http.StatusBadRequest)))
	}
	if id, err = c.GetInt(":id"); err != nil || id < 0 {
		if err == nil {
			err = fmt.Errorf("error: El id debe ser positivo")
		}
		panic(e.Error(funcion, err, strconv.Itoa(http.StatusBadRequest)))
	}
	if detailed, err = c.GetBool("detailed", false); err != nil {
		panic(e.Error(funcion, err, strconv.Itoa(http.StatusBadRequest)))
	}

	if v, err := transaccionmovimientos.Get(idType, id, detailed); err != nil {
		panic(err)
	} else {
		c.Data["json"] = v
		c.Ctx.Output.SetStatus(200)
	}
	c.ServeJSON()
}
