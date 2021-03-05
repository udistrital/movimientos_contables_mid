package models

import "time"

type TransaccionMovimientos struct {
	ConsecutivoId    int
	Etiquetas        string
	Descripcion      string
	FechaTransaccion time.Time
	Activo           bool
	Movimientos      []struct {
		TerceroId        *int
		CuentaId         string
		NombreCuenta     string
		TipoMovimientoId int
		Valor            float64
		Descripcion      string
		Activo           bool
	} `json:"movimientos"`
}
