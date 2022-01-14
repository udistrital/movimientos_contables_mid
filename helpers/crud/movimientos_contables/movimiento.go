package movimientos_contables

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	e "github.com/udistrital/utils_oas/errorctrl"
	r "github.com/udistrital/utils_oas/request"

	"github.com/udistrital/movimientos_contables_mid/helpers"
)

// GetMovimientos retorna las transacciones segun los criterios
func GetMovimientos(query string, fields []string, limit int, offset int, movimientos interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetMovimientos"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))

	var fullResponse map[string]interface{}
	params := url.Values{}
	params.Add("query", query)
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))
	url := beego.AppConfig.String("MovimientosContablesCrudService") + "/movimiento?" + params.Encode()
	if resp, err := r.GetJsonTest(url, &fullResponse); err != nil || resp.StatusCode != http.StatusOK {
		status := http.StatusBadGateway
		if err == nil { // resp.StatusCode != http.StatusOK
			err = fmt.Errorf("undesired status code - %s", http.StatusText(resp.StatusCode))
			status = resp.StatusCode
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - r.GetJsonTest(url, &fullResponse)", err, strconv.Itoa(status))
	}

	helpers.LimpiezaRespuestaRefactor(fullResponse, &movimientos)
	return
}

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
