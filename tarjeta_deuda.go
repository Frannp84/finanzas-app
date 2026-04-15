package main

import "fmt"

func verDeudaTarjeta() {

	if len(tarjetas) == 0 {

		fmt.Println("No hay tarjetas cargadas")
		return

	}

	indiceTarjeta := elegirTarjeta()

	t := tarjetas[indiceTarjeta]

	deuda := t.LimiteTotal - t.LimiteDisponible

	fmt.Println("\n--- DEUDA TARJETA ---")

	fmt.Println("\nTarjeta:", t.Nombre)

	fmt.Println("\nLimite total:",
		formatearMoneda(t.LimiteTotal))

	fmt.Println("Disponible:",
		formatearMoneda(t.LimiteDisponible))

	fmt.Println("\nDeuda actual:",
		formatearMoneda(deuda))

}