package models

import "time"

type TransaccionMovimientos struct {
	ConsecutivoId    int
	Etiquetas        string
	Descripcion      string
	FechaTransaccion time.Time
	Activo           bool
	Movimientos      []MovimientoResumido   `json:"movimientos"`
	Comprobante      map[string]interface{} `json:"Comprobante"`
}
