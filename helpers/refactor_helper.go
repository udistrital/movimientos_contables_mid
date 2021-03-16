package helpers

import (
	"encoding/json"
	_ "fmt"
)

func LimpiezaRespuestaRefactor(respuesta map[string]interface{}, v interface{}) {

	b, err := json.Marshal(respuesta["Data"])
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, v)
}

func LimpiezaRespuestaRefactorBody(respuesta map[string]interface{}, v interface{}) {

	b, err := json.Marshal(respuesta["Body"])
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, v)
}
