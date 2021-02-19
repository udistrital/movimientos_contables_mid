package models

type Consecutivo struct {
	Id                int     `orm:"column(id);pk;auto"`
	ContextoId        int     `orm:"column(contexto_id)"`
	Year              float64 `orm:"column(year)"`
	Consecutivo       int     `orm:"column(consecutivo)"`
	Descripcion       string  `orm:"column(descripcion);null"`
	Activo            bool    `orm:"column(activo)"`
	FechaCreacion     string  `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion string  `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
}
