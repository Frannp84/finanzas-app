package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var tarjetas []Tarjeta

func guardarTarjetas() {

	file, _ := os.Create("tarjetas.json")
	defer file.Close()

	json.NewEncoder(file).Encode(tarjetas)

}

func cargarTarjetas() {

	file, err := os.Open("tarjetas.json")

	if err != nil {
		return
	}

	defer file.Close()

	json.NewDecoder(file).Decode(&tarjetas)

}
func agregarTarjeta() {

	var nombre string
	var limite float64
	var cierre int
	var vencimiento int

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nNombre de la tarjeta:")
	nombre, _ = reader.ReadString('\n')

	nombre = strings.TrimSpace(nombre)

	fmt.Println("\nLimite total:")
	fmt.Scanln(&limite)

	fmt.Println("\nDia de cierre:")
	fmt.Scanln(&cierre)

	fmt.Println("\nDia de vencimiento:")
	fmt.Scanln(&vencimiento)

	fmt.Println("\nCuenta desde la que se paga:")

	indiceCuenta := elegirCuenta()

	nueva := Tarjeta{

		Nombre: nombre,

		LimiteTotal:      limite,
		LimiteDisponible: limite,

		Cierre:      cierre,
		Vencimiento: vencimiento,

		CuentaAsociada: indiceCuenta,
	}

	tarjetas = append(tarjetas, nueva)

	guardarTarjetas()

	fmt.Println("Tarjeta creada correctamente")

}
func verTarjetas() {

	fmt.Println("\n--- TARJETAS ---")

	for i, t := range tarjetas {

		fmt.Printf("%d - %s\n", i+1, t.Nombre)

		fmt.Println("Limite total:", formatearMoneda(t.LimiteTotal))

		fmt.Println("Disponible:", formatearMoneda(t.LimiteDisponible))

		fmt.Println("Cierre:", t.Cierre)

		fmt.Println("Vencimiento:", t.Vencimiento)

		fmt.Println("Cuenta asociada:",
			cuentas[t.CuentaAsociada].Nombre)

		fmt.Println()

	}

}
func elegirTarjeta() int {

	fmt.Println("\nSeleccione tarjeta:")

	for i, t := range tarjetas {

		fmt.Printf("%d %s\n",
			i+1,
			t.Nombre)

	}

	var opcion int

	fmt.Print("\nNumero: ")
	fmt.Scanln(&opcion)

	if opcion <= 0 || opcion > len(tarjetas) {

		fmt.Println("Tarjeta invalida")
		return elegirTarjeta()

	}

	return opcion - 1

}
func pagarTarjeta() {

	if len(tarjetas) == 0 {

		fmt.Println("No hay tarjetas cargadas")
		return

	}

	indiceTarjeta := elegirTarjeta()

	var monto float64

	fmt.Println("\nMonto pagado:")
	fmt.Scanln(&monto)

	fmt.Println("\nCuenta desde donde se paga:")

	indiceCuenta := elegirCuenta()

	// bajar saldo de cuenta
	cuentas[indiceCuenta].Saldo -= monto

	// recuperar limite disponible
	tarjetas[indiceTarjeta].LimiteDisponible += monto

	// evitar que supere el limite total
	if tarjetas[indiceTarjeta].LimiteDisponible > tarjetas[indiceTarjeta].LimiteTotal {

		tarjetas[indiceTarjeta].LimiteDisponible = tarjetas[indiceTarjeta].LimiteTotal

	}

	guardarCuentas()
	guardarTarjetas()

	fmt.Println("\nPago registrado correctamente")

}
