package helpers

import "testing"

func TestMovimientosBalanceadosTruncaTotalesAMilesimas(t *testing.T) {
	tests := []struct {
		nombre   string
		debito   float64
		credito  float64
		esperado bool
	}{
		{
			nombre:   "acepta totales iguales despues de sumar y truncar",
			debito:   905965.0047 + 14824070.011099998,
			credito:  15730035.015799996,
			esperado: true,
		},
		{
			nombre:   "rechaza diferencia que permanece despues de truncar",
			debito:   10.1239,
			credito:  10.124,
			esperado: false,
		},
		{
			nombre:   "trunca en lugar de redondear",
			debito:   1.9999,
			credito:  1.999,
			esperado: true,
		},
		{
			nombre:   "trunca valores negativos hacia cero",
			debito:   -1.9999,
			credito:  -1.999,
			esperado: true,
		},
	}

	for _, test := range tests {
		t.Run(test.nombre, func(t *testing.T) {
			if obtenido := movimientosBalanceados(test.debito, test.credito); obtenido != test.esperado {
				t.Fatalf("movimientosBalanceados(%v, %v) = %v; esperado %v", test.debito, test.credito, obtenido, test.esperado)
			}
		})
	}
}

func TestTruncarAMilesimasConDecimalExacto(t *testing.T) {
	const valor = 15730035.015
	const esperado int64 = 15730035015

	if obtenido := truncarAMilesimas(valor); obtenido != esperado {
		t.Fatalf("truncarAMilesimas(%v) = %d; esperado %d", valor, obtenido, esperado)
	}
}
