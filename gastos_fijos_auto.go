package main

func generarMovimientosGastosFijos(mes int, anio int) {

	for _, g := range gastosFijos {

		if anio < g.AnioInicio {
			continue
		}

		if anio == g.AnioInicio && mes < g.MesInicio {
			continue
		}

		existe := false

		for _, m := range movimientos {

			if m.Detalle == g.Nombre &&
				m.Mes == mes &&
				m.Anio == anio {

				existe = true
				break
			}

		}

		if existe {
			continue
		}

		mov := Movimiento{

			Tipo: "gasto",

			Monto: g.Monto,

			Detalle: g.Nombre,

			Categoria: g.Categoria,

			Anio: anio,
			Mes:  mes,
			Dia:  g.Dia,

			Cuotas: 1,

			MedioPago: g.MedioPago,

			Cuenta:  g.Cuenta,
			Tarjeta: g.Tarjeta,
		}

		movimientos = append(movimientos, mov)

		// impacto saldo inmediato
		if g.MedioPago == "Cuenta" {

			cuentas[g.Cuenta].Saldo -= g.Monto

		}

		if g.MedioPago == "Tarjeta" {

			tarjetas[g.Tarjeta].LimiteDisponible -= g.Monto

		}

	}

	guardarDatos()
	guardarCuentas()
	guardarTarjetas()

}
