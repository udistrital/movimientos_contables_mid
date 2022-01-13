package transaccionmovimientos

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	e "github.com/udistrital/utils_oas/errorctrl"

	"github.com/udistrital/movimientos_contables_mid/helpers/crud/movimientos_contables"
)

func Get(tipoDeId string, id int, conMovimientos bool) (transaccion map[string]interface{}, outputError map[string]interface{}) {
	const funcion string = "Get"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))

	vars := map[string]interface{}{
		"tipoDeId":       tipoDeId,
		"id":             id,
		"conMovimientos": conMovimientos,
	}
	logs.Debug(vars)

	criterios := map[string]string{
		"consecutivo": "ConsecutivoId",
		"transaccion": "Id",
	}

	var query string
	if k, ok := criterios[tipoDeId]; !ok {
		err := fmt.Errorf("criterio '%s' no valido", tipoDeId)
		outputError = e.Error(funcion+" - criterios[tipoDeId]", err, strconv.Itoa(http.StatusBadRequest))
		return
	} else {
		query = fmt.Sprintf("%s:%d", k, id)
	}

	var transacciones []map[string]interface{}
	if err := movimientos_contables.GetTransaccionesByQuery(query, &transacciones); err != nil {
		outputError = err
		return
	}
	if len(transacciones) == 0 {
		err := fmt.Errorf("no existe transaccion con '%s': %d", tipoDeId, id)
		outputError = e.Error(funcion+" - len(transacciones) == 0", err, strconv.Itoa(http.StatusNotFound))
		return
	} else if len(transacciones) > 1 {
		err := fmt.Errorf("hay mas de una transaccion con '%s': %d", tipoDeId, id)
		outputError = e.Error(funcion+" - len(transacciones) > 1", err, strconv.Itoa(http.StatusConflict))
		return
	}

	transaccion = transacciones[0]
	if conMovimientos {
		query := fmt.Sprintf("TransaccionId:%v", transaccion["Id"])
		fields := []string{"Activo", "CuentaId", "Descripcion", "Id", "NombreCuenta", "TerceroId", "TipoMovimientoId", "Valor"}
		var movimientos []map[string]interface{}
		if err := movimientos_contables.GetMovimientos(query, fields, -1, 0, &movimientos); err != nil {
			outputError = err
		}
		if len(movimientos) > 0 {
			transaccion["movimientos"] = movimientos
		} else {
			transaccion["movimientos"] = []interface{}{}
		}
	}

	return
}
