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
	var concepto models.Concepto
	var valor_debito float64 = 0
	var valor_credito float64 = 0
	var etiquetas models.Etiquetas
	var error_valor_movimiento bool = false

	//validaciones consecutivo
	if v.ConsecutivoId > 0 {
		if response, err := getJsonTest(beego.AppConfig.String("ConsecutivosCrudService")+"/consecutivo/?query=Id:"+v.ConsecutivoId, &respuesta_peticion); (err == nil) && (response == 200) {
			if len(respuesta_peticion["Data"].([]interface{})[0].(map[string]interface{})) == 0 {
				v.ErrorTransaccion = "Error: el consecutivo ingresado no se encuentra en la base de datos \n"
			}
			if response, err := getJsonTest(beego.AppConfig.String("MovimientosContablesCrudService")+"/transaccion/?query=ConsecutivoId:"+v.ConsecutivoId, &respuesta_peticion); (err == nil) && (response == 200) {
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

	etiquetas = models.Etiquetas{}
	b, errEtiqueta := json.Marshal(v.Etiquetas)
	if errEtiqueta != nil {
		panic(errEtiqueta)
	}
	json.Unmarshal(b, etiquetas)
	if response, err := getJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/concepto/"+etiquetas.ConceptoId, &respuesta_peticion); (err == nil) && (response == 200) {
		if len(respuesta_peticion["Body"].([]interface{})[0].(map[string]interface{})) == 0 {
			v.ErrorTransaccion += "Error : el Id de concepto ingresado: " + etiquetas.ConceptoId + " no se encuentra registrado \n"
		} else {
			concepto = models.Concepto{}
			LimpiezaRespuestaRefactor(respuesta_peticion, &concepto)
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos4", "err": err.Error(), "status": "404"}
		return outputError
	}

	if time.Now().Before(v.FechaTransaccion) {
		v.ErrorTransaccion += "Error: la fecha ingresada es mayor a la fecha actual \n"
	}

	//validacion movimientos

	for _, movimiento := range v.Movimientos {
		// validacion de existencia de cuenta
		nodo_cuenta_contable = models.NodoCuentaContable{}
		if response, err := getJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/nodo_cuenta_contable/"+movimiento.CuentaId, &respuesta_peticion); (err == nil) && (response == 200) {
			if len(respuesta_peticion["Body"].([]interface{})[0].(map[string]interface{})) == 0 {
				v.ErrorTransaccion += "Error: el Numero de cuenta ingresado: " + movimiento.CuentaId + " no se encuentra registrado \n"
			} else {
				LimpiezaRespuestaRefactor(respuesta_peticion, &nodo_cuenta_contable)
			}
		} else { //If nodo_cuenta_contable get
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos5", "err": err.Error(), "status": "404"}
			return outputError
		}
		if movimiento.TerceroId != nil {
			if !nodo_cuenta_contable.RequiereTercero {
				v.ErrorTransaccion += "Error: la cuenta: " + movimiento.CuentaId + " no deberia tener registrado un tercero "
			}
			if response, err := getJsonTest(beego.AppConfig.String("TercerosCrudService")+"/tercero/"+movimiento.TerceroId, &respuesta_peticion); (err == nil) && (response == 200) {
				if len(respuesta_peticion_tercero) == 0 { //probar funcionamiento de esta validacion
					v.ErrorTransaccion += "Error: el tercero ingresado con id: " + movimiento.TerceroId + ", para el numero de cuenta: " + movimiento.CuentaId + " no se encuentra registrado \n"
				}
			} else { //If tercero get
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos6", "err": err.Error(), "status": "404"}
				return outputError
			}
		} else {
			if nodo_cuenta_contable.RequiereTercero {
				v.ErrorTransaccion += "Error: la cuenta: " + movimiento.CuentaId + " requiere que se registre un tercero"
			}
		}

		if concepto != (models.Concepto{}) {
			if movimiento.TipoMovimientoId != 1 && movimiento.TipoMovimientoId != 2 { // pendiente validar el id de tipomoviento en la tabla parametro
				v.ErrorTransaccion += "Error: el tipo de movimiento registrado para la cuenta: " + movimiento.CuentaId + " no es valido"
				error_valor_movimiento = true
			} else {
				if movimiento.TipoMovimientoId == 1 {
					if concepto.CuentaDebito == movimiento.CuentaId {
						if movimiento.Valor < 0 {
							v.ErrorTransaccion += "Error: el numero de cuenta: " + movimiento.CuentaId + "  registra un valor invalido \n"
							error_valor_movimiento = true
						} else {
							valor_debito += movimiento.Valor
						}
					} else {
						v.ErrorTransaccion += "Error: el numero de cuenta: " + movimiento.CuentaId + "  registra un valor en la columna incorrecta \n"
					}
				} else {
					if concepto.CuentaCredito == movimiento.CuentaId {
						if movimiento.Valor < 0 {
							v.ErrorTransaccion += "Error: el numero de cuenta: " + movimiento.CuentaId + "  registra un valor invalido \n"
							error_valor_movimiento = true
						} else {
							valor_credito += movimiento.Valor
						}
					} else {
						v.ErrorTransaccion += "Error: el numero de cuenta: " + movimiento.CuentaId + "  registra un valor en la columna incorrecta \n"
						error_valor_movimiento = true
					}
				}
			}
		}
	}

	// validacion sumas iguales
	if error_valor_movimiento {
		v.ErrorTransaccion += "Error: no es posible verificar las sumas iguales ya que hay errores en los valores de los movimientos \n"
	} else if valor_debito != valor_credito {
		v.ErrorTransaccion += "Error: los movimientos no cumplen con el requerimiento de sumas iguales \n"
	}

	//definicion del estado de una transaccion
	if v.ErrorTransaccion != "" {
		v.EstadoId = 2 //pendiente definir id estsado transaccion
	} else {
		v.EstadoId = 1 //pendiente definir id estsado transaccion
	}

	if err := sendJson(beego.AppConfig.String("MovimientosContablesCrudService")+"/transaccion", "POST", &response, transaccion); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos", "err": err.Error(), "status": "502"}
		return outputError
	}

	return nil
}
