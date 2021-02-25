package helpers

import (
	_ "fmt"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/movimientos_contables_mid/models"
)

func RegistroTransaccionMovimientos(v models.TransaccionMovimientos) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var response interface{}
	var respuesta_peticion map[string]interface{}

	//validaciones consecutivo
	if v.ConsecutivoId > 0 {
		if response, err := getJsonTest(beego.AppConfig.String("ProtocolCrudConsecutivos")+"://"+beego.AppConfig.String("UrlCrudConsecutivos")+"/"+beego.AppConfig.String("NsCrudConsecutivos")+"/consecutivo/?query=Id:"+v.ConsecutivoId, &respuesta_peticion); (err == nil) && (response == 200) {
			if len(respuesta_peticion["Data"].([]interface{})[0].(map[string]interface{})) == 0 {
				v.ErrorTransaccion = "Error: el consecutivo ingresado no se encuentra en la base de datos \n"
			}
			if response, err := getJsonTest(beego.AppConfig.String("ProtocolCrudMovimientosContables")+"://"+beego.AppConfig.String("UrlCrudMovimientosContables")+"/"+beego.AppConfig.String("NsCrudMovimientosContables")+"/transaccion/?query=ConsecutivoId:"+v.ConsecutivoId, &respuesta_peticion); (err == nil) && (response == 200) {
				if len(respuesta_peticion["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
					v.ErrorTransaccion += "Error: el consecutivo ingresado ya tiene una transaccion asociada \n"
				}
			} else { //If transaccion get
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos2", "err": err.Error(), "status": "404"}
				return outputError
			}
		} else { //If consecutivo get
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos3", "err": err.Error(), "status": "404"}
			return outputError
		}
	} else {
		v.ConsecutivoId = -1 //pendiente por arreglar
		v.ErrorTransaccion = "Error: no se ha ingresado un consecutivo valido \n"
	}

	if time.Now().Before(v.FechaTransaccion) {
		v.ErrorTransaccion += "Error: la fecha ingresada es mayor a la fecha actual \n"
	}

	if v.ErrorTransaccion != "" {
		v.EstadoId = 2 //pendiente definir id estsado transaccion
	} else {
		v.EstadoId = 1 //pendiente definir id estsado transaccion
	}

	if err := sendJson(beego.AppConfig.String("ProtocolCrudMovimientosContables")+"://"+beego.AppConfig.String("UrlCrudMovimientosContables")+"/"+beego.AppConfig.String("NsCrudMovimientosContables")+"/transaccion", "POST", &response, transaccion); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos", "err": err.Error(), "status": "502"}
		return outputError
	}

	return nil
}
