package models

type NodoCuentaContable struct {
	Codigo              string
	Hijos               []string
	Padre               *string
	Nombre              string
	Nivel               int
	DetalleCuentaID     string
	NaturalezaCuentaID  string
	CodigoCuentaAlterna string
	Ajustable           bool
	MonedaID            string
	RequiereTercero     bool
	CentroDecostosID    string
	Nmnc                bool
	FechaCreacion       string
	FechaModificacion   string
	Activo              bool
}
