package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:TransaccionMovimientosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:TransaccionMovimientosController"],
        beego.ControllerComments{
            Method: "PostTransaccionMovimientos",
            Router: "/transaccion_movimientos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
