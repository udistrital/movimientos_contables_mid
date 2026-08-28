package movimientos_contables

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	e "github.com/udistrital/utils_oas/errorctrl"
	r "github.com/udistrital/utils_oas/request"

	"github.com/udistrital/movimientos_contables_mid/helpers"
)

// GetTransaccionesByQuery retorna las transacciones buscando por query
func GetTransaccionesByQuery(ctx context.Context, query string, transacciones interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetTransaccionesByQuery"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))

	var fullResponse map[string]interface{}
	url := beego.AppConfig.String("MovimientosContablesCrudService") + "/transaccion?query=" + url.QueryEscape(query)
	if status, err := r.GetWithContext(ctx, url, &fullResponse); err != nil {
		if status == 0 {
			status = http.StatusBadGateway
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - request.GetWithContext(ctx, url, &fullResponse)", err, strconv.Itoa(status))
		return
	}

	helpers.LimpiezaRespuestaRefactor(fullResponse, &transacciones)
	return
}

func PostTransaccion(in interface{}, out interface{}) (outputError map[string]interface{}) {
	const funcion string = "PostTransaccion"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))

	url := beego.AppConfig.String("MovimientosContablesCrudService") + "/transaccion"
	if err := r.SendJson(url, "POST", &out, in); err != nil {
		logs.Error(err)
		status := strconv.Itoa(http.StatusBadGateway)
		outputError = e.Error(funcion+" - r.SendJson(url, \"POST\", &out, in)", err, status)
	}
	return
}
