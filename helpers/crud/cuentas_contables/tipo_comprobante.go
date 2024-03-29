package cuentas_contables

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	e "github.com/udistrital/utils_oas/errorctrl"
	r "github.com/udistrital/utils_oas/request"
)

func GetTipoComprobanteById(id string, tipo_comprobante interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetTipoComprobanteById"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))

	url := beego.AppConfig.String("CuentasContablesCrudService") + "/tipo_comprobante/" + id
	if resp, err := r.GetJsonTest(url, &tipo_comprobante); err != nil || resp.StatusCode != http.StatusOK {
		status := http.StatusBadGateway
		if err == nil { // resp.StatusCode != http.StatusOK
			err = fmt.Errorf("undesired status code - %s", http.StatusText(resp.StatusCode))
			status = resp.StatusCode
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - r.GetJsonTest(url, &tipo_comprobante)", err, strconv.Itoa(status))
	}
	return
}
