package movimientos_contables

import (
	"fmt"
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
func GetTransaccionesByQuery(query string, transacciones []map[string]interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetTransaccionesByQuery"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))

	var fullResponse map[string]interface{}
	url := beego.AppConfig.String("MovimientosContablesCrudService") + "/transaccion?query=" + url.QueryEscape(query)
	if resp, err := r.GetJsonTest(url, &fullResponse); err != nil || resp.StatusCode != http.StatusOK {
		status := http.StatusBadGateway
		if err == nil { // resp.StatusCode != http.StatusOK
			err = fmt.Errorf("undesired status code - %s", http.StatusText(resp.StatusCode))
			status = resp.StatusCode
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - r.GetJsonTest(url, &fullResponse)", err, strconv.Itoa(status))
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
