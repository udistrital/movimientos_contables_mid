package models

import "time"

type Transaccion struct {
	Id                int
	ConsecutivoId     int
	Etiquetas         string
	Descripcion       string
	ErrorTransaccion  string
	EstadoId          int
	FechaTransaccion  time.Time
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}
