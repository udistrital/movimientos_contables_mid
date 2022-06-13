package cuentas_contables

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/movimientos_contables_mid/helpers"
	e "github.com/udistrital/utils_oas/errorctrl"
	r "github.com/udistrital/utils_oas/request"
)

func GetNodoCuentaContableByCuentaId(cuentaId string, nodo interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetNodoCuentaContableByCuentaId"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))
	var fullResponse map[string]interface{}
	url := beego.AppConfig.String("CuentasContablesCrudService") + "/nodo_cuenta_contable/" + cuentaId
	if resp, err := r.GetJsonTest(url, &fullResponse); err != nil || resp.StatusCode != http.StatusOK {
		status := http.StatusBadGateway
		if err == nil { // resp.StatusCode != http.StatusOK
			err = fmt.Errorf("undesired status code - %s", http.StatusText(resp.StatusCode))
			status = resp.StatusCode
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - r.GetJsonTest(url, &fullResponse)", err, strconv.Itoa(status))
	}
	helpers.LimpiezaRespuestaRefactorBody(fullResponse, &nodo)
	return
}

func GetNodoCuentaContableWorker(id string, c chan interface{}) {
	var nodo interface{}
	outputError := GetNodoCuentaContableByCuentaId(id, &nodo)
	if outputError != nil {
		logs.Warn(outputError)
		c <- nil
	} else {
		c <- nodo
	}
}

func DeleteNodoCuentaContableByCuentaId(cuentaId string) (outputError map[string]interface{}) {
	const funcion string = "DeleteNodoCuentaContableByCuentaId"
	var exist interface{}
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))
	var fullResponseExist map[string]interface{}
	var fullResponseDelete map[string]interface{}
	url := beego.AppConfig.String("MovimientosContablesCrudService") + "/movimiento/exist/" + cuentaId

	if resp, err := r.GetJsonTest(url, &fullResponseExist); err == nil && resp.StatusCode == http.StatusOK {

		helpers.LimpiezaRespuestaRefactor(fullResponseExist, &exist)

		if exist.(bool) {

			outputError = e.Error(funcion+" - exist", err, strconv.Itoa(http.StatusBadRequest))
			return
		}

		url = beego.AppConfig.String("CuentasContablesCrudService") + "/nodo_cuenta_contable/" + cuentaId

		if errd := r.SendJson(url, "DELETE", &fullResponseDelete, nil); errd != nil || fullResponseDelete["Code"].(float64) != 200 {

			status := http.StatusBadGateway
			if err == nil {
				err = fmt.Errorf("undesired status code - %s", http.StatusText(resp.StatusCode))
				status = resp.StatusCode
			}
			logs.Error(errd)
			outputError = e.Error(funcion+" - r.DeleteJsonTest(url, &fullResponse)", err, strconv.Itoa(status))
		}

	} else {

		logs.Error(err)
		outputError = e.Error(funcion+" - r.GetJsonTest(url, &fullResponse)", err, strconv.Itoa(http.StatusBadGateway))
	}

	return
}
