// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/movimientos_contables_mid/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/transaccion_movimientos",
			beego.NSInclude(
				&controllers.TransaccionMovimientosController{},
			),
		),
		beego.NSNamespace("/movimientos",
			beego.NSInclude(
				&controllers.MovimientosController{},
			),
		),
		beego.NSNamespace("/cuenta_contable",
			beego.NSInclude(
				&controllers.CuentaContableController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
