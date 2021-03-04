package helpers

import (
	"encoding/json"
	"fmt"
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
			outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var response map[string]interface{}
	var respuesta_peticion map[string]interface{}
	var nodo_cuenta_contable models.NodoCuentaContable
	var valor_debito float64 = 0
	var valor_credito float64 = 0
	var etiquetas models.Etiquetas
	var error_valor_movimiento bool = false
	var transaccion models.Transaccion

	//validaciones consecutivo
	if v.ConsecutivoId > 0 { //el consecutivo se encuentra registrado dentro del cuerpo de la peticion
		//verificacion en la consulta del consecutivo
		if response, err := getJsonTest(beego.AppConfig.String("ConsecutivosCrudService")+"/consecutivo/"+strconv.Itoa(v.ConsecutivoId), &respuesta_peticion); (err == nil) && (response == 200 || response == 404) {
			//verificacion del contenido de la consulta del consecutivo
			if response == 404 {
				//v.ErrorTransaccion = "Error: el consecutivo ingresado no se encuentra en la base de datos \n"
				v.ErrorTransaccion = "E1,"
			}
			//consulta de una transaccion asociada al consecutivo
			if response, err := getJsonTest(beego.AppConfig.String("MovimientosContablesCrudService")+"/transaccion/?query=ConsecutivoId:"+strconv.Itoa(v.ConsecutivoId), &respuesta_peticion); (err == nil) && (response == 200) {
				//verificacion del contenido de la consulta de la transaccion
				if len(respuesta_peticion["Data"].([]interface{})[0].(map[string]interface{})) != 0 {
					//v.ErrorTransaccion += "Error: el consecutivo ingresado ya tiene una transaccion asociada \n"
					v.ErrorTransaccion += "E2,"
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
		v.ConsecutivoId = 0
		//v.ErrorTransaccion = "Error: no se ha ingresado un consecutivo valido \n"
		v.ErrorTransaccion = "E3,"
	}

	// validaciones etiquetas
	etiquetas = models.Etiquetas{}
	//verificacion de error en el json.marshall
	errEtiqueta := json.Unmarshal([]byte(v.Etiquetas), &etiquetas)
	if errEtiqueta != nil {
		panic(errEtiqueta.Error() + "error en etiquetas")
	}

	fmt.Println(etiquetas)
	fmt.Println(etiquetas.ComprobanteId)
	fmt.Println(etiquetas.TipoComprobanteId)
	//validacion tipo comprobante
	if etiquetas.TipoComprobanteId != "" {
		fmt.Println(respuesta_peticion)
		if response, err := getJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/tipo_comprobante/"+etiquetas.TipoComprobanteId, &respuesta_peticion); (err == nil) && (response == 200) {
			fmt.Println(respuesta_peticion)
			if (respuesta_peticion["Type"].(interface{}).(string)) == "error" {
				//v.ErrorTransaccion += "Error: el Id de tipo de comprobante ingresado no se encuentra registrado \n"
				v.ErrorTransaccion += "E4,"
				fmt.Println(v.ErrorTransaccion)
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos4", "err": err.Error(), "status": "404"}
			return outputError
		}
	}
	//validacion comprobante
	fmt.Println(respuesta_peticion)
	if etiquetas.ComprobanteId != "" {
		if response, err := getJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/comprobante/"+etiquetas.ComprobanteId, &respuesta_peticion); (err == nil) && (response == 200) {
			fmt.Println(respuesta_peticion)
			if (respuesta_peticion["Type"].(interface{}).(string)) == "error" {
				//v.ErrorTransaccion += "Error: el Id de comprobante ingresado no se encuentra registrado \n"
				v.ErrorTransaccion += "E5,"
				fmt.Println(v.ErrorTransaccion)
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos5", "err": err.Error(), "status": "404"}
			return outputError
		}
	}
	//validacion de fecha de la transaccion
	fmt.Println("fecha transaccion")
	if time.Now().Before(v.FechaTransaccion) {
		//v.ErrorTransaccion += "Error: la fecha ingresada es mayor a la fecha actual \n"
		v.ErrorTransaccion += "E6,"
	}
	fmt.Println("movimientos")
	if v.Movimientos != nil {

		//validacion movimientos
		fmt.Println("array movimientos")
		for _, movimiento := range v.Movimientos {
			// validacion de existencia de cuenta
			fmt.Println("for movimientos")
			nodo_cuenta_contable = models.NodoCuentaContable{}
			//consulta de la cuenta asociada al movimiento
			if response, err := getJsonTest(beego.AppConfig.String("CuentasContablesCrudService")+"/nodo_cuenta_contable/"+movimiento.CuentaId, &respuesta_peticion); (err == nil) && (response == 200) {
				//verificacion del contenido de la consulta de la cuenta
				if (respuesta_peticion["Body"].(interface{}).(map[string]interface{})["Codigo"]) == "" { //revisar a detalle
					//v.ErrorTransaccion += "Error: el Numero de cuenta ingresado: " + movimiento.CuentaId + " no se encuentra registrado \n"
					v.ErrorTransaccion += "E7,"
				} else { //mapeo de la respuesta en el objeto nodo_cuenta_contable
					LimpiezaRespuestaRefactor(respuesta_peticion, &nodo_cuenta_contable)
				}
			} else { //If nodo_cuenta_contable get
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos6", "err": err.Error(), "status": "404"}
				return outputError
			}
			//verificacion de la presencia del tercero en el cuerpo de la peticion
			if movimiento.TerceroId != nil {
				//comprobar si la cuenta asociada al movimiento requiere registrar un tercero
				if !nodo_cuenta_contable.RequiereTercero {
					//v.ErrorTransaccion += "Error: la cuenta: " + movimiento.CuentaId + " no deberia tener registrado un tercero "
					v.ErrorTransaccion += "E8,"
				}
				//verificar si el tercero existe
				if response, err := getJsonTest(beego.AppConfig.String("TercerosCrudService")+"/tercero/"+strconv.Itoa(*movimiento.TerceroId), &respuesta_peticion); (err == nil) && (response == 200 || response == 404) {
					//verificacion del contenido de la consulta del tercero
					if response == 404 { //probar funcionamiento de esta validacion
						//v.ErrorTransaccion += "Error: el tercero ingresado con id: " + strconv.Itoa(*movimiento.TerceroId) + ", para el numero de cuenta: " + movimiento.CuentaId + " no se encuentra registrado \n"
						v.ErrorTransaccion += "E9,"
					}
				} else { //If tercero get
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos7", "err": err.Error(), "status": "404"}
					return outputError
				}
			} else {
				//comprobar si la cuenta asociada al movimiento requiere un tercero
				if nodo_cuenta_contable.RequiereTercero {
					//v.ErrorTransaccion += "Error: la cuenta: " + movimiento.CuentaId + " requiere que se registre un tercero"
					v.ErrorTransaccion += "E10,"
				}
			}

			//verificar que el id del tipo de movimiento coincida con la informacion registrada en parametros
			if movimiento.TipoMovimientoId != 1 && movimiento.TipoMovimientoId != 2 { // pendiente validar el id de tipomoviento en la tabla parametro
				//v.ErrorTransaccion += "Error: el tipo de movimiento registrado para la cuenta: " + movimiento.CuentaId + " no es valido"
				v.ErrorTransaccion += "E11,"
				error_valor_movimiento = true
			} else { //id de tipo movimiento invalido
				if movimiento.Valor < 0 {
					//v.ErrorTransaccion += "Error: el numero de cuenta: " + movimiento.CuentaId + "  registra un valor invalido \n"
					v.ErrorTransaccion += "E12,"
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
	} else {
		fmt.Println("error no hay movimientos")
		//v.ErrorTransaccion += "Error: no se ha ingresado ningun movimiento dentro de la transaccion \n"
		v.ErrorTransaccion += "E13,"
	}

	// validacion valores del movimiento
	fmt.Println("error valor movimiento")
	if error_valor_movimiento {
		//v.ErrorTransaccion += "Error: no es posible verificar las sumas iguales ya que hay errores en los valores de los movimientos \n"
		v.ErrorTransaccion += "E14,"
	} else if valor_debito != valor_credito { // validacion de sumas iguales
		//v.ErrorTransaccion += "Error: los movimientos no cumplen con el requerimiento de sumas iguales \n"
		v.ErrorTransaccion += "E15,"
	}

	//definicion del estado de una transaccion
	fmt.Println("error transaccion ")
	if v.ErrorTransaccion != "" {
		v.EstadoId = 2 //pendiente definir id estsado transaccion
		outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos8", "err": v.ErrorTransaccion, "status": "400"}
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

	fmt.Println("envio transaccion ")
	a, _ := json.Marshal(transaccion)
	fmt.Println(a)
	fmt.Println("------------------------------")
	fmt.Println(transaccion.Etiquetas)
	if err := sendJson(beego.AppConfig.String("MovimientosContablesCrudService")+"/transaccion", "POST", &response, transaccion); err != nil {
		fmt.Println(response)
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos9", "err": err.Error(), "status": "502"}
		return outputError
	}
	fmt.Println(response)
	fmt.Println("envio movimientos")
	if v.Movimientos != nil {
		if len(response["Data"].(interface{}).(map[string]interface{})) != 0 {
			LimpiezaRespuestaRefactor(response, &transaccion)
			for _, movimiento := range v.Movimientos {
				movimiento.TransaccionId = transaccion.Id
				if err := sendJson(beego.AppConfig.String("MovimientosContablesCrudService")+"/movimiento", "POST", &response, movimiento); err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/RegistroTransaccionMovimientos10", "err": err.Error(), "status": "502"}
					return outputError
				}
			}
		}
	}
	return outputError
}
