package helpers

import (
	_ "fmt"

	"github.com/udistrital/utils_oas/formatdata"
)

func LimpiezaRespuestaRefactor(respuesta map[string]interface{}, v interface{}) {
	formatdata.FillStruct(respuesta["Data"], &v)
}

func LimpiezaRespuestaRefactorBody(respuesta map[string]interface{}, v interface{}) {
	formatdata.FillStruct(respuesta["Body"], &v)
}
