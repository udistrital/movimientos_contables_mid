package models

type Parametro struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
	FechaCreacion     string
	FechaModificacion string
	TipoParametroId   int
	ParametroPadreId  int
}
