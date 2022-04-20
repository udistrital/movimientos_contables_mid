package transaccionmovimientos

import (
	"fmt"
	"net/http"
	"strconv"

	e "github.com/udistrital/utils_oas/errorctrl"

	"github.com/udistrital/movimientos_contables_mid/helpers/crud/cuentas_contables"
	"github.com/udistrital/movimientos_contables_mid/helpers/crud/movimientos_contables"
)

func Get(tipoDeId string, id int, conMovimientos bool) (transaccion map[string]interface{}, outputError map[string]interface{}) {
	const funcion string = "Get"

	query, err := validarCriterios(tipoDeId, id)
	if err != nil {
		outputError = e.Error(funcion+" - len(transacciones) == 0", err, strconv.Itoa(http.StatusNotFound))
		return
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
	compchan := make(chan interface{})
	movchan := make(chan interface{})
	etiquetaString := fmt.Sprintf("%v", transaccion["Etiquetas"])
	transaccionId := fmt.Sprintf("%v", transaccion["Id"])
	go cuentas_contables.GetComprobanteWorker(etiquetaString, compchan)
	go movimientos_contables.GetMovimientosWorker(transaccionId, conMovimientos, movchan)
	transaccion["Comprobante"] = <-compchan
	transaccion["Movimientos"] = <-movchan

	return
}

func validarCriterios(tipoDeId string, id int) (query string, err error) {
	const funcion string = "Get.validarCriterios"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))

	criterios := map[string]string{
		"consecutivo": "ConsecutivoId",
		"transaccion": "Id",
	}

	if k, ok := criterios[tipoDeId]; !ok {
		err = fmt.Errorf("criterio '%s' no valido", tipoDeId)
	} else {
		query = fmt.Sprintf("%s:%d", k, id)
	}
	return
}
