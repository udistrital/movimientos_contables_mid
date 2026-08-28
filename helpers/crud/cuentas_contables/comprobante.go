package cuentas_contables

import (
	"context"
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

func GetComprobanteById(ctx context.Context, id string, comprobante interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetComprobanteById"
	var fullResponse map[string]interface{}
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))

	url := beego.AppConfig.String("CuentasContablesCrudService") + "/comprobante/" + id
	if status, err := r.GetWithContext(ctx, url, &fullResponse); err != nil {
		if status == 0 {
			status = http.StatusBadGateway
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - request.GetWithContext(ctx, url, &fullResponse)", err, strconv.Itoa(status))
		return
	}
	helpers.LimpiezaRespuestaRefactorBody(fullResponse, &comprobante)
	return
}

func GetComprobanteWorker(ctx context.Context, etiquetaString string, c chan interface{}) {
	var etiqueta map[string]string
	if err := json.Unmarshal([]byte(etiquetaString), &etiqueta); err == nil && etiqueta["ComprobanteId"] != "" {
		id := fmt.Sprintf("%v", etiqueta["ComprobanteId"])
		var comprobante interface{}
		outputError := GetComprobanteById(ctx, id, &comprobante)
		if outputError != nil {
			logs.Warn(outputError)
			c <- nil
		} else {
			c <- comprobante
		}
	} else {
		c <- nil
	}

}
