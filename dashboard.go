package main

import (
	"fmt"
	"time"
)

func verDashboard() {

	now := time.Now()

	mes := int(now.Month())
	anio := now.Year()

	// genera automaticamente gastos fijos del mes
	generarMovimientosGastosFijos(mes, anio)

	mostrarDashboard(mes, anio)

}

func verDashboardPeriodo() {

	var mes int
	var anio int

	fmt.Println("\nAño:")
	fmt.Scanln(&anio)

	fmt.Println("\nMes:")

	for i := 1; i <= 12; i++ {

		fmt.Printf("%d - %s\n",
			i,
			nombreMes(i))

	}

	fmt.Print("\nNumero: ")
	fmt.Scanln(&mes)

	generarMovimientosGastosFijos(mes, anio)

	mostrarDashboard(mes, anio)

}

func mostrarDashboard(mes int, anio int) {

	var totalARS float64
	var totalUSD float64

	for _, c := range cuentas {

		switch c.Moneda {

		case "ARS":

			totalARS += c.Saldo

		case "USD":

			totalUSD += c.Saldo

		}

	}

	var totalIngresos float64
	var totalGastos float64

	for _, m := range movimientos {

		valorCuota := m.Monto / float64(m.Cuotas)

		for i := 0; i < m.Cuotas; i++ {

			mesCuota := m.Mes + i
			anioCuota := m.Anio

			for mesCuota > 12 {

				mesCuota -= 12
				anioCuota++

			}

			if mesCuota == mes && anioCuota == anio {

				if m.Tipo == "ingreso" {

					totalIngresos += valorCuota

				} else {

					totalGastos += valorCuota

				}

			}

		}

	}

	resultado := totalIngresos - totalGastos

	var porcentajeAhorro float64

	if totalIngresos > 0 {

		porcentajeAhorro = (resultado / totalIngresos) * 100

	}

	var cuotasFuturas float64

	for _, m := range movimientos {

		if m.Tipo != "gasto" {
			continue
		}

		if m.Cuotas <= 1 {
			continue
		}

		valorCuota := m.Monto / float64(m.Cuotas)

		for i := 1; i < m.Cuotas; i++ {

			cuotasFuturas += valorCuota

		}

	}

	var deudaTarjetas float64

	for _, t := range tarjetas {

		deudaTarjetas += (t.LimiteTotal - t.LimiteDisponible)

	}

	fmt.Println("\n--- DASHBOARD ---")

	fmt.Printf("\nPeriodo: %s %d\n",
		nombreMes(mes),
		anio)

	fmt.Println("\nPatrimonio ARS:",
		formatearMoneda(totalARS))

	fmt.Printf("Patrimonio USD: %.2f\n",
		totalUSD)

	fmt.Println("\nIngresos del mes:",
		formatearMoneda(totalIngresos))

	fmt.Println("Gastos del mes:",
		formatearMoneda(totalGastos))

	fmt.Println("\nResultado del mes:",
		formatearMoneda(resultado))

	fmt.Printf("Ahorro: %.1f%%\n",
		porcentajeAhorro)

	fmt.Println("\nCuotas futuras:",
		formatearMoneda(cuotasFuturas))

	fmt.Println("Deuda tarjetas:",
		formatearMoneda(deudaTarjetas))

}
