package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:CuentaContableController"] = append(beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:CuentaContableController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:MovimientosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:MovimientosController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:TransaccionMovimientosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:TransaccionMovimientosController"],
        beego.ControllerComments{
            Method: "PostTransaccionMovimientos",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:TransaccionMovimientosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:TransaccionMovimientosController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/:idType/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:TransaccionMovimientosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/movimientos_contables_mid/controllers:TransaccionMovimientosController"],
        beego.ControllerComments{
            Method: "PostTransaccionMovimientosDeprecada",
            Router: "/transaccion_movimientos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
