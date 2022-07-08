package terceros

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	e "github.com/udistrital/utils_oas/errorctrl"
	r "github.com/udistrital/utils_oas/request"
)

func GetTerceroById(terceroId string, tercero interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetTerceroById"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))
	url := beego.AppConfig.String("TercerosCrudService") + "/tercero/" + terceroId
	if resp, err := r.GetJsonTest(url, &tercero); err != nil || resp.StatusCode != http.StatusOK {
		status := http.StatusBadGateway
		if err == nil { // resp.StatusCode != http.StatusOK
			err = fmt.Errorf("undesired status code - %s", http.StatusText(resp.StatusCode))
			status = resp.StatusCode
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - r.GetJsonTest(url, &tercero)", err, strconv.Itoa(status))
	}
	return
}

func GetTerceroWorker(id *int, c chan interface{}) {
	var tercero interface{}
	if id != nil {
		outputError := GetTerceroById(strconv.Itoa(*id), &tercero)
		if outputError != nil {
			logs.Warn(outputError)
			c <- nil
		} else {
			c <- tercero
		}
	} else {
		c <- nil
	}

}
