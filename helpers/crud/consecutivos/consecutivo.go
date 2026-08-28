package consecutivos

import (
	"context"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/movimientos_contables_mid/helpers"
	e "github.com/udistrital/utils_oas/errorctrl"
	r "github.com/udistrital/utils_oas/request"
)

func GetById(ctx context.Context, id int, consecutivo interface{}) (outputError map[string]interface{}) {
	const funcion string = "GetById"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", strconv.Itoa(http.StatusInternalServerError))
	var fullResponse map[string]interface{}
	url := beego.AppConfig.String("ConsecutivosCrudService") + "/consecutivo/" + strconv.Itoa(id)
	if status, err := r.GetWithContext(ctx, url, &fullResponse); err != nil {
		if status == 0 {
			status = http.StatusBadGateway
		}
		logs.Error(err)
		outputError = e.Error(funcion+" - request.GetWithContext(ctx, url, &fullResponse)", err, strconv.Itoa(status))
		return
	}
	helpers.LimpiezaRespuestaRefactor(fullResponse, &consecutivo)
	return
}

func GetConsecutivoWorker(ctx context.Context, id int, c chan interface{}) {
	var consecutivo interface{}
	outputError := GetById(ctx, id, &consecutivo)
	if outputError != nil {
		logs.Warn(outputError)
		c <- nil
	} else {
		c <- consecutivo
	}
}
