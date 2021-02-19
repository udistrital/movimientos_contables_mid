package models

import "time"

type Transaccion struct {
	Id                int       `orm:"column(id);pk;auto"`
	ConsecutivoId     int       `orm:"column(consecutivo_id)"`
	Etiquetas         string    `orm:"column(etiquetas);type(json);null"`
	Descripcion       string    `orm:"column(descripcion);null"`
	ErrorTransaccion  string    `orm:"column(error_transaccion);null"`
	EstadoId          int       `orm:"column(estado_id)"`
	FechaTransaccion  time.Time `orm:"column(fecha_transaccion);type(timestamp without time zone)"`
	Activo            bool      `orm:"column(activo)"`
	FechaCreacion     string    `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion string    `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
}
