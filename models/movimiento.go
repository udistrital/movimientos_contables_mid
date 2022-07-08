package models

type Movimiento struct {
	Id                int
	TerceroId         *int
	CuentaId          string
	NombreCuenta      string
	TipoMovimientoId  int
	Valor             float64
	Descripcion       string
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
	TransaccionId     *Transaccion
	SaldoAnterior     float64
	NuevoSaldo        float64
	Cuenta            interface{}
	Tercero           interface{}
	Consecutivo       interface{}
	Comprobante       interface{}
}

type MovimientoResumido struct {
	TerceroId        *int
	CuentaId         string
	NombreCuenta     string
	TipoMovimientoId int
	Valor            float64
	Descripcion      string
	Activo           bool
}
