package consecutivos

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

func GetById(id int, consecutivo interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetById"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))
	var fullResponse map[string]interface{}
	url := beego.AppConfig.String("ConsecutivosCrudService") + "/consecutivo/" + strconv.Itoa(id)
	if resp, err := r.GetJsonTest(url, &fullResponse); err != nil || resp.StatusCode != http.StatusOK {
		status := http.StatusBadGateway
		if err == nil { // resp.StatusCode != http.StatusOK
			err = fmt.Errorf("undesired status code - %s", http.StatusText(resp.StatusCode))
			status = resp.StatusCode
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - r.GetJsonTest(url, &consecutivo)", err, strconv.Itoa(status))
	}
	helpers.LimpiezaRespuestaRefactor(fullResponse, &consecutivo)
	return
}

func GetConsecutivoWorker(id int, c chan interface{}) {
	var consecutivo interface{}
	outputError := GetById(id, &consecutivo)
	if outputError != nil {
		logs.Warn(outputError)
		c <- nil
	} else {
		c <- consecutivo
	}
}
