package models

type TransaccionMovimientos struct {
	*Transaccion
	Movimientos []*Movimiento
}
