package models

type Movimiento struct {
	Id                int          `orm:"column(id);pk;auto"`
	TerceroId         int          `orm:"column(tercero_id);null"`
	CuentaId          string       `orm:"column(cuenta_id)"`
	NombreCuenta      string       `orm:"column(nombre_cuenta)"`
	TipoMovimientoId  int          `orm:"column(tipo_movimiento_id)"`
	Valor             float64      `orm:"column(valor)"`
	Descripcion       string       `orm:"column(descripcion);null"`
	Activo            bool         `orm:"column(activo)"`
	FechaCreacion     string       `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion string       `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
	TransaccionId     *Transaccion `orm:"column(transaccion_id);rel(fk)"`
}
