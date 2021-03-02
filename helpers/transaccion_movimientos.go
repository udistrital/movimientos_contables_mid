package helpers

import (
	"encoding/json"
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
	var respuesta_peticion_tercero []interface{}
	var nodo_cuenta_contable models.NodoCuentaContable
	var valor_debito float64 = 0
	var valor_credito float64 = 0
	var etiquetas models.Etiquetas
	var error_valor_movimiento bool = false
	var transaccion models.Transaccion

	//validaciones consecutivo
	if v.ConsecutivoId > 0 { //el consecutivo se encuentra registrado dentro del cuerpo de la peticion
		//verificacion en la consulta del consecutivo
		if response, err := getJsonTest(beego.AppConfig.String("ConsecutivosCrudService")+"/consecutivo/?query=Id:"+v.ConsecutivoId, &respuesta_peticion); (err == nil) && (response == 200) {
			//verificacion del contenido de la consulta del consecutivo
			if len(respuesta_peticion["Data"].([]interface{})[0].(map[string]interface{})) == 0 {
				v.ErrorTransaccion = "Error: el consecutivo ingresado no se encuentra en la base de datos \n"
			}
			//consulta de una transaccion asociada al consecutivo
			if response, err := getJsonTest(beego.AppConfig.String("MovimientosContablesCrudService")+"/transaccion/?query=ConsecutivoId:"+v.ConsecutivoId, &respuesta_peticion); (err == nil) && (response == 200) {
				//verificacion del contenido de la consulta de la transaccion
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
	} else { //consecutivo invalido
		v.ConsecutivoId = -1 //pendiente por arreglar
		v.ErrorTransaccion = "Error: no se ha ingresado un consecutivo valido \n"
	}

	etiquetas = models.Etiquetas{}
	b, errEtiqueta := json.Marshal(v.Etiquetas)
	//verificacion de error en el json.marshall
	if errEtiqueta != nil {
		panic(errEtiqueta)
	}
	json.Unmarshal(b, etiquetas)

	if response, err := getJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/tipo_comprobante/"+etiquetas.TipoComprobanteId, &respuesta_peticion); (err == nil) && (response == 200) {
		if len(respuesta_peticion["Body"].([]interface{})[0].(map[string]interface{})) == 0 {
		}
	}

	//validacion de fecha de la transaccion
	if time.Now().Before(v.FechaTransaccion) {
		v.ErrorTransaccion += "Error: la fecha ingresada es mayor a la fecha actual \n"
	}

	//validacion movimientos
	for _, movimiento := range v.Movimientos {
		// validacion de existencia de cuenta
		nodo_cuenta_contable = models.NodoCuentaContable{}
		//consulta de la cuenta asociada al movimiento
		if response, err := getJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/nodo_cuenta_contable/"+movimiento.CuentaId, &respuesta_peticion); (err == nil) && (response == 200) {
			//verificacion del contenido de la consulta de la cuenta
			if len(respuesta_peticion["Body"].([]interface{})[0].(map[string]interface{})) == 0 {
				v.ErrorTransaccion += "Error: el Numero de cuenta ingresado: " + movimiento.CuentaId + " no se encuentra registrado \n"
			} else { //mapeo de la respuesta en el objeto nodo_cuenta_contable
				LimpiezaRespuestaRefactor(respuesta_peticion, &nodo_cuenta_contable)
			}
		} else { //If nodo_cuenta_contable get
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos4", "err": err.Error(), "status": "404"}
			return outputError
		}
		//verificacion de la presencia del tercero en el cuerpo de la peticion
		if movimiento.TerceroId != nil {
			//comprobar si la cuenta asociada al movimiento requiere registrar un tercero
			if !nodo_cuenta_contable.RequiereTercero {
				v.ErrorTransaccion += "Error: la cuenta: " + movimiento.CuentaId + " no deberia tener registrado un tercero "
			}
			//verificar si el tercero existe
			if response, err := getJsonTest(beego.AppConfig.String("TercerosCrudService")+"/tercero/"+movimiento.TerceroId, &respuesta_peticion); (err == nil) && (response == 200) {
				//verificacion del contenido de la consulta del tercero
				if len(respuesta_peticion_tercero) == 0 { //probar funcionamiento de esta validacion
					v.ErrorTransaccion += "Error: el tercero ingresado con id: " + movimiento.TerceroId + ", para el numero de cuenta: " + movimiento.CuentaId + " no se encuentra registrado \n"
				}
			} else { //If tercero get
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos5", "err": err.Error(), "status": "404"}
				return outputError
			}
		} else {
			//comprobar si la cuenta asociada al movimiento requiere un tercero
			if nodo_cuenta_contable.RequiereTercero {
				v.ErrorTransaccion += "Error: la cuenta: " + movimiento.CuentaId + " requiere que se registre un tercero"
			}
		}

		//verificar que el id del tipo de movimiento coincida con la informacion registrada en parametros
		if movimiento.TipoMovimientoId != 1 && movimiento.TipoMovimientoId != 2 { // pendiente validar el id de tipomoviento en la tabla parametro
			v.ErrorTransaccion += "Error: el tipo de movimiento registrado para la cuenta: " + movimiento.CuentaId + " no es valido"
			error_valor_movimiento = true
		} else { //id de tipo movimiento invalido
			if movimiento.Valor < 0 {
				v.ErrorTransaccion += "Error: el numero de cuenta: " + movimiento.CuentaId + "  registra un valor invalido \n"
				error_valor_movimiento = true
			} else {
				// tipo de movimiento = debito
				if movimiento.TipoMovimientoId == 1 {
					valor_debito += movimiento.Valor
				} else { //tipo de movimiento = credito
					valor_credito += movimiento.Valor
				}
			}
		}
	}

	// validacion valores del movimiento
	if error_valor_movimiento {
		v.ErrorTransaccion += "Error: no es posible verificar las sumas iguales ya que hay errores en los valores de los movimientos \n"
	} else if valor_debito != valor_credito { // validacion de sumas iguales
		v.ErrorTransaccion += "Error: los movimientos no cumplen con el requerimiento de sumas iguales \n"
	}

	//definicion del estado de una transaccion
	if v.ErrorTransaccion != "" {
		v.EstadoId = 2 //pendiente definir id estsado transaccion
		outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos6", "err": v.ErrorTransaccion, "status": "400"}
	} else {
		v.EstadoId = 1 //pendiente definir id estsado transaccion
	}

	//
	transaccion = models.Transaccion{}
	transaccion.FechaTransaccion = v.FechaTransaccion
	transaccion.Etiquetas = v.Etiquetas
	transaccion.EstadoId = v.EstadoId
	transaccion.ErrorTransaccion = v.ErrorTransaccion
	transaccion.ConsecutivoId = v.ConsecutivoId
	transaccion.Activo = true
	transaccion.Descripcion = v.Descripcion
	transaccion.ConsecutivoId = v.ConsecutivoId

	if err := sendJson(beego.AppConfig.String("MovimientosContablesCrudService")+"/transaccion", "POST", &response, transaccion); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos7", "err": err.Error(), "status": "502"}
		return outputError
	}

	for _, movimiento := range v.Movimientos {
		if err := sendJson(beego.AppConfig.String("MovimientosContablesCrudService")+"/movimiento", "POST", &response, movimiento); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos8", "err": err.Error(), "status": "502"}
			return outputError
		}
	}

	return outputError
}
