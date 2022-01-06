package movimientos_contables

import (
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	e "github.com/udistrital/utils_oas/errorctrl"
	r "github.com/udistrital/utils_oas/request"
)

func PostMovimiento(in interface{}, out interface{}) (outputError map[string]interface{}) {
	const funcion string = "PostMovimiento"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))

	url := beego.AppConfig.String("MovimientosContablesCrudService") + "/movimiento"
	if err := r.SendJson(url, "POST", &out, in); err != nil {
		logs.Error(err)
		status := strconv.Itoa(http.StatusBadGateway)
		outputError = e.Error(funcion+" - r.SendJson(url, \"POST\", &out, in)", err, status)
	}
	return
}
