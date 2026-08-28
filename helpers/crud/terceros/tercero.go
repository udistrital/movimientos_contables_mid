package terceros

import (
	"context"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	e "github.com/udistrital/utils_oas/errorctrl"
	r "github.com/udistrital/utils_oas/request"
)

func GetTerceroById(ctx context.Context, terceroId string, tercero interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetTerceroById"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))
	url := beego.AppConfig.String("TercerosCrudService") + "/tercero/" + terceroId
	if status, err := r.GetWithContext(ctx, url, &tercero); err != nil {
		if status == 0 {
			status = http.StatusBadGateway
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - request.GetWithContext(ctx, url, &tercero)", err, strconv.Itoa(status))
	}
	return
}

func GetTerceroWorker(ctx context.Context, id *int, c chan interface{}) {
	var tercero interface{}
	if id != nil {
		outputError := GetTerceroById(ctx, strconv.Itoa(*id), &tercero)
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
