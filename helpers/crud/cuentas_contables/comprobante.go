package cuentas_contables

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/movimientos_contables_mid/helpers"
	e "github.com/udistrital/utils_oas/errorctrl"
	r "github.com/udistrital/utils_oas/request"
)

func GetComprobanteById(id string, comprobante interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetComprobanteById"
	var fullResponse map[string]interface{}
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))

	url := beego.AppConfig.String("CuentasContablesCrudService") + "/comprobante/" + id
	if resp, err := r.GetJsonTest(url, &fullResponse); err != nil || resp.StatusCode != http.StatusOK {
		status := http.StatusBadGateway
		if err == nil { // resp.StatusCode != http.StatusOK
			err = fmt.Errorf("undesired status code - %s", http.StatusText(resp.StatusCode))
			status = resp.StatusCode
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - r.GetJsonTest(url, &fullResponse)", err, strconv.Itoa(status))
	}
	helpers.LimpiezaRespuestaRefactorBody(fullResponse, &comprobante)
	return
}

func GetComprobanteWorker(etiquetaString string, c chan interface{}) {
	var etiqueta map[string]string
	if err := json.Unmarshal([]byte(etiquetaString), &etiqueta); err == nil {
		id := fmt.Sprintf("%v", etiqueta["ComprobanteId"])
		var comprobante interface{}
		outputError := GetComprobanteById(id, &comprobante)
		if outputError != nil {
			c <- nil
		} else {
			c <- comprobante
		}
	} else {
		c <- nil
	}

}
