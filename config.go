package main

import "fmt"

func menuConfiguracion() {

	var opcion int

	for {

		fmt.Println("\n--- CONFIGURACION ---")

		fmt.Println("\nCATEGORIAS")
		fmt.Println("1 Ver categorias")
		fmt.Println("2 Agregar categoria")
		fmt.Println("3 Eliminar categoria")

		fmt.Println("\nDATOS")
		fmt.Println("4 Reset movimientos")
		fmt.Println("5 Reset presupuestos")

		fmt.Println("\nCUENTAS")
		fmt.Println("6 Ver cuentas")
		fmt.Println("7 Agregar cuenta")

		fmt.Println("\nTARJETAS")
		fmt.Println("8 Ver Tarjetas")
		fmt.Println("9 Agregar Tarjeta")
		fmt.Println("10 Pagar Tarjeta")

		fmt.Println("\nGASTOS FIJOS")
		fmt.Println("11 Ver gastos fijos")
		fmt.Println("12 Agregar gasto fijo")

		fmt.Println("\n0 volver")

		fmt.Print("\nOpcion: ")
		fmt.Scanln(&opcion)

		switch opcion {

		case 1:
			verCategorias()

		case 2:
			agregarCategoria()

		case 3:
			eliminarCategoria()

		case 4:
			borrarMovimientosSeguro()

		case 5:
			borrarPresupuestosSeguro()

		case 6:
			verCuentas()

		case 7:
			agregarCuenta()

		case 8:
			verTarjetas()

		case 9:
			agregarTarjeta()

		case 10:
			pagarTarjeta()

		case 11:
    		verGastosFijos()

		case 12:
    		agregarGastoFijo()

		case 0:
			return

		default:
			fmt.Println("Opcion invalida")

		}

	}

}
