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
