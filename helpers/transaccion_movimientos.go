package helpers

import (
	"encoding/json"
	_ "fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/movimientos_contables_mid/models"
)

func RegistroTransaccionMovimientos(v models.TransaccionMovimientos) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var response map[string]interface{}
	var respuesta_peticion map[string]interface{}
	var nodo_cuenta_contable models.NodoCuentaContable
	var valor_debito float64 = 0
	var valor_credito float64 = 0
	var etiquetas models.Etiquetas
	var movimiento_envio models.Movimiento
	var error_valor_movimiento bool = false
	var transaccion models.Transaccion = models.Transaccion{}

	transaccion.FechaTransaccion = v.FechaTransaccion
	transaccion.Etiquetas = v.Etiquetas
	transaccion.ConsecutivoId = v.ConsecutivoId
	transaccion.Activo = v.Activo
	transaccion.Descripcion = v.Descripcion
	transaccion.ConsecutivoId = v.ConsecutivoId

	//validaciones consecutivo
	transaccion.ErrorTransaccion = ""
	if transaccion.ConsecutivoId > 0 { //el consecutivo se encuentra registrado dentro del cuerpo de la peticion
		//verificacion en la consulta del consecutivo
		if response, err := getJsonTest(beego.AppConfig.String("ConsecutivosCrudService")+"/consecutivo/"+strconv.Itoa(transaccion.ConsecutivoId), &respuesta_peticion); (err == nil) && (response == 200 || response == 404) {
			//verificacion del contenido de la consulta del consecutivo
			if response == 404 {
				transaccion.ErrorTransaccion = "Error: el consecutivo ingresado no se encuentra en la base de datos \n"
			}
			//consulta de una transaccion asociada al consecutivo
			if response, err := getJsonTest(beego.AppConfig.String("MovimientosContablesCrudService")+"/transaccion/?query=ConsecutivoId:"+strconv.Itoa(transaccion.ConsecutivoId), &respuesta_peticion); (err == nil) && (response == 200) {
				//verificacion del contenido de la consulta de la transaccion
				if len(respuesta_peticion["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
					transaccion.ErrorTransaccion += "Error: el consecutivo ingresado ya tiene una transaccion asociada \n"
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
	} else { //consecutivo invalido
		transaccion.ConsecutivoId = 0
		transaccion.ErrorTransaccion = "Error: no se ha ingresado un consecutivo valido \n"
	}

	// validaciones etiquetas
	if transaccion.Etiquetas != "" {
		etiquetas = models.Etiquetas{}
		//verificacion de error en el json.marshall
		errEtiqueta := json.Unmarshal([]byte(transaccion.Etiquetas), &etiquetas)
		if errEtiqueta != nil {
			panic(errEtiqueta.Error() + "error en etiquetas")
		}
		//validacion tipo comprobante
		if etiquetas.TipoComprobanteId != "" {
			if response, err := getJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/tipo_comprobante/"+etiquetas.TipoComprobanteId, &respuesta_peticion); (err == nil) && (response == 200) {
				if (respuesta_peticion["Code"].(float64)) == 500 {
					transaccion.ErrorTransaccion += "Error: el Id de tipo de comprobante ingresado no se encuentra registrado \n"
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos4", "err": err.Error(), "status": "404"}
				return outputError
			}
		}
		//validacion comprobante
		if etiquetas.ComprobanteId != "" {
			if response, err := getJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/comprobante/"+etiquetas.ComprobanteId, &respuesta_peticion); (err == nil) && (response == 200) {
				if (respuesta_peticion["Code"].(float64)) == 500 {
					transaccion.ErrorTransaccion += "Error: el Id de comprobante ingresado no se encuentra registrado \n"
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos5", "err": err.Error(), "status": "404"}
				return outputError
			}
		}
	}
	//validacion de fecha de la transaccion
	if time.Now().Before(transaccion.FechaTransaccion) {
		transaccion.ErrorTransaccion += "Error: la fecha ingresada es mayor a la fecha actual \n"
	}
	if v.Movimientos != nil {
		//validacion movimientos
		for i, movimiento := range v.Movimientos {
			// validacion de existencia de cuenta
			nodo_cuenta_contable = models.NodoCuentaContable{}
			//consulta de la cuenta asociada al movimiento
			if response, err := getJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/nodo_cuenta_contable/"+movimiento.CuentaId, &respuesta_peticion); (err == nil) && (response == 200) {
				//verificacion del contenido de la consulta de la cuenta
				if (respuesta_peticion["Body"].(map[string]interface{})["Codigo"]) == "" { //revisar a detalle
					transaccion.ErrorTransaccion += "Error: el Numero de cuenta ingresado: " + movimiento.CuentaId + " no se encuentra registrado \n"
				} else { //mapeo de la respuesta en el objeto nodo_cuenta_contable
					LimpiezaRespuestaRefactorBody(respuesta_peticion, &nodo_cuenta_contable)
				}
			} else { //If nodo_cuenta_contable get
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos6", "err": err.Error(), "status": "404"}
				return outputError
			}
			if nodo_cuenta_contable.Codigo != "" {
				v.Movimientos[i].NombreCuenta = nodo_cuenta_contable.Nombre
			}
			//verificacion de la presencia del tercero en el cuerpo de la peticion
			if movimiento.TerceroId != nil {
				//comprobar si la cuenta asociada al movimiento requiere registrar un tercero
				if nodo_cuenta_contable.Codigo != "" {
					if !nodo_cuenta_contable.RequiereTercero {
						transaccion.ErrorTransaccion += "Error: la cuenta: " + movimiento.CuentaId + " no deberia tener registrado un tercero "
					}
				}
				//verificar si el tercero existe
				if response, err := getJsonTest(beego.AppConfig.String("TercerosCrudService")+"/tercero/"+strconv.Itoa(*movimiento.TerceroId), &respuesta_peticion); (err == nil) && (response == 200 || response == 404) {
					//verificacion del contenido de la consulta del tercero
					if response == 404 { //probar funcionamiento de esta validacion
						transaccion.ErrorTransaccion += "Error: el tercero ingresado con id: " + strconv.Itoa(*movimiento.TerceroId) + ", para el numero de cuenta: " + movimiento.CuentaId + " no se encuentra registrado \n"
					}
				} else { //If tercero get
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos7", "err": err.Error(), "status": "404"}
					return outputError
				}
			} else {
				//comprobar si la cuenta asociada al movimiento requiere un tercero
				if nodo_cuenta_contable.Codigo != "" {
					if nodo_cuenta_contable.RequiereTercero {
						transaccion.ErrorTransaccion += "Error: la cuenta: " + movimiento.CuentaId + " requiere que se registre un tercero"
					}
				}
			}
			//verificar que el id del tipo de movimiento coincida con la informacion registrada en parametros
			if movimiento.TipoMovimientoId != 344 && movimiento.TipoMovimientoId != 345 { // pendiente validar el id de tipomoviento en la tabla parametro
				transaccion.ErrorTransaccion += "Error: el tipo de movimiento registrado para la cuenta: " + movimiento.CuentaId + " no es valido"
				error_valor_movimiento = true
			} else { //id de tipo movimiento invalido
				if movimiento.Valor < 0 {
					transaccion.ErrorTransaccion += "Error: el numero de cuenta: " + movimiento.CuentaId + "  registra un valor invalido \n"
					error_valor_movimiento = true
				} else {
					// tipo de movimiento = debito
					if movimiento.TipoMovimientoId == 344 {
						valor_debito += movimiento.Valor
					} else { //tipo de movimiento = credito
						valor_credito += movimiento.Valor
					}
				}
			}
		}
	} else {
		transaccion.ErrorTransaccion += "Error: no se ha ingresado ningun movimiento dentro de la transaccion \n"
	}
	// validacion valores del movimiento
	if error_valor_movimiento {
		transaccion.ErrorTransaccion += "Error: no es posible verificar las sumas iguales ya que hay errores en los valores de los movimientos \n"
	} else if valor_debito != valor_credito { // validacion de sumas iguales
		transaccion.ErrorTransaccion += "Error: los movimientos no cumplen con el requerimiento de sumas iguales \n"
	}
	//definicion del estado de una transaccion
	if transaccion.ErrorTransaccion != "" {
		transaccion.EstadoId = 343 //id de la tabla parametro (parametros crud) que relaciona una transaccion invalida
		outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos8", "err": transaccion.ErrorTransaccion, "status": "400"}
	} else {
		transaccion.EstadoId = 342 //id de la tabla parametro (parametros crud) que relaciona una transaccion valida
	}
	if err := sendJson(beego.AppConfig.String("MovimientosContablesCrudService")+"/transaccion", "POST", &response, transaccion); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos9", "err": err.Error(), "status": "502"}
		return outputError
	}
	if v.Movimientos != nil {
		if len(response["Data"].(map[string]interface{})) != 0 {
			LimpiezaRespuestaRefactor(response, &transaccion)
			for _, movimiento := range v.Movimientos {
				movimiento_envio = models.Movimiento{}
				movimiento_envio.TerceroId = movimiento.TerceroId
				movimiento_envio.CuentaId = movimiento.CuentaId
				movimiento_envio.NombreCuenta = movimiento.NombreCuenta
				movimiento_envio.TipoMovimientoId = movimiento.TipoMovimientoId
				movimiento_envio.Valor = movimiento.Valor
				movimiento_envio.Descripcion = movimiento.Descripcion
				movimiento_envio.Activo = movimiento.Activo
				movimiento_envio.TransaccionId = &transaccion

				if err := sendJson(beego.AppConfig.String("MovimientosContablesCrudService")+"/movimiento", "POST", &response, movimiento_envio); err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos10", "err": err.Error(), "status": "502"}
					return outputError
				}
			}
		}
	}
	return outputError
}
