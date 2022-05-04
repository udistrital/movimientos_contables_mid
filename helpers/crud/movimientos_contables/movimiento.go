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
	"github.com/udistrital/movimientos_contables_mid/helpers/crud/consecutivos"
	"github.com/udistrital/movimientos_contables_mid/helpers/crud/cuentas_contables"
	"github.com/udistrital/movimientos_contables_mid/helpers/crud/terceros"
	"github.com/udistrital/movimientos_contables_mid/models"
)

// GetMovimientos retorna las transacciones segun los criterios
func GetMovimientos(query string, fields []string, limit int, offset int, sortby []string, order []string, detailfields []string, m interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetMovimientos"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))
	var movimientos []models.Movimiento
	var fullResponse map[string]interface{}
	params := url.Values{}
	params.Add("query", query)
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}
	if len(sortby) > 0 {
		params.Add("sortby", strings.Join(sortby, ","))
	}
	if len(order) > 0 {
		params.Add("order", strings.Join(order, ","))
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
		GetMovimientoDetalle(&movimiento, detailfields)
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
		detailfields := []string{
			"Cuenta",
			"Tercero",
		}
		var movimientos interface{}
		outputError := GetMovimientos(query, fields, -1, 0, nil, nil, detailfields, &movimientos)
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

func GetMovimientoDetalle(movimiento *models.Movimiento, fields []string) {
	nodochan := make(chan interface{})
	terchan := make(chan interface{})
	compchan := make(chan interface{})
	conschan := make(chan interface{})
	tercero, cuenta, consecutivo, comprobante, all := false, false, false, false, false
	if len(fields) > 0 {
		for _, field := range fields {
			switch field {
			case "Tercero":
				tercero = true
			case "Cuenta":
				cuenta = true
			case "Consecutivo":
				consecutivo = true
			case "Comprobante":
				comprobante = true
			}
		}
	} else {
		all = true
	}
	if tercero || all {
		go terceros.GetTerceroWorker(movimiento.TerceroId, terchan)
	} else {
		close(terchan)
	}
	if (cuenta || all) && movimiento.CuentaId != "" {
		go cuentas_contables.GetNodoCuentaContableWorker(movimiento.CuentaId, nodochan)
	} else {
		close(nodochan)
	}
	if (consecutivo || all) && movimiento.TransaccionId != nil {
		go consecutivos.GetConsecutivoWorker(movimiento.TransaccionId.ConsecutivoId, conschan)
	} else {
		close(conschan)
	}
	if (comprobante || all) && movimiento.TransaccionId != nil {
		go cuentas_contables.GetComprobanteWorker(movimiento.TransaccionId.Etiquetas, compchan)
	} else {
		close(compchan)
	}
	movimiento.Cuenta = <-nodochan
	movimiento.Tercero = <-terchan
	movimiento.Consecutivo = <-conschan
	movimiento.Comprobante = <-compchan
}
