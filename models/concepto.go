package models

type Concepto struct {
	*General      `bson:"inline"`
	ID            string
	Nombre        string
	CuentaDebito  string
	CuentaCredito string
	SistemaID     string
}
