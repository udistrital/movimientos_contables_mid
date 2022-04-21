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
	f "github.com/udistrital/utils_oas/formatdata"
	r "github.com/udistrital/utils_oas/request"

	"github.com/udistrital/movimientos_contables_mid/helpers"
	"github.com/udistrital/movimientos_contables_mid/models"
)

// GetMovimientos retorna las transacciones segun los criterios
func GetMovimientos(query string, fields []string, limit int, offset int, m interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetMovimientos"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))
	var movimientos []models.Movimiento
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
	for i, movimiento := range movimientos {
		var respuesta_peticion map[string]interface{}
		var cuenta interface{}
		if _, err := r.GetJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/nodo_cuenta_contable/"+movimiento.CuentaId, &respuesta_peticion); err == nil {
			helpers.LimpiezaRespuestaRefactorBody(respuesta_peticion, &cuenta)
			fmt.Print(cuenta)
		} else {
			logs.Error(err)
			cuenta = nil
		}
		movimiento.Cuenta = &cuenta
		movimientos[i] = movimiento
	}
	f.FillStruct(movimientos, &m)
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

func GetMovimientosWorker(id string, conMovimientos bool, c chan interface{}) {
	if conMovimientos {
		query := fmt.Sprintf("TransaccionId:%v", id)
		fields := []string{
			"Activo",
			"CuentaId",
			"Descripcion",
			"Id",
			"NombreCuenta",
			"TerceroId",
			"TipoMovimientoId",
			"Valor",
		}
		var movimientos interface{}
		outputError := GetMovimientos(query, fields, -1, 0, &movimientos)
		if outputError != nil {
			logs.Warn(outputError)
			c <- nil
		} else {
			c <- movimientos
		}
	} else {
		c <- nil
	}

}
