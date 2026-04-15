package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type GastoFijo struct {

	Nombre string

	Monto float64

	Categoria string

	Dia int

	MedioPago string

	Cuenta int
	Tarjeta int

	AnioInicio int
	MesInicio int

}

var gastosFijos []GastoFijo

func guardarGastosFijos() {

	file, _ := os.Create("gastos_fijos.json")
	defer file.Close()

	json.NewEncoder(file).Encode(gastosFijos)

}

func cargarGastosFijos() {

	file, err := os.Open("gastos_fijos.json")

	if err != nil {
		return
	}

	defer file.Close()

	json.NewDecoder(file).Decode(&gastosFijos)

}

func agregarGastoFijo() {

	var nombre string
	var monto float64
	var dia int
	var anio int
	var mes int

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nNombre del gasto fijo:")
	nombre, _ = reader.ReadString('\n')
	nombre = strings.TrimSpace(nombre)

	fmt.Println("\nMonto:")
	fmt.Scanln(&monto)

	categoria := elegirCategoria("gasto")

	fmt.Println("\nDia del gasto:")
	fmt.Scanln(&dia)

	fmt.Println("\nDesde año:")
	fmt.Scanln(&anio)

	fmt.Println("\nDesde mes:")
	fmt.Scanln(&mes)

	fmt.Println("\nMedio de pago:")
	fmt.Println("1 Cuenta / Efectivo")
	fmt.Println("2 Tarjeta credito")

	var medio int
	fmt.Scanln(&medio)

	var cuenta int
	var tarjeta int = -1
	var medioTexto string

	if medio == 2 {

		medioTexto = "Tarjeta"

		fmt.Println("\nSeleccione tarjeta:")
		tarjeta = elegirTarjeta()

	} else {

		medioTexto = "Cuenta"

		fmt.Println("\nSeleccione cuenta:")
		cuenta = elegirCuenta()

	}

	g := GastoFijo{

		Nombre: nombre,
		Monto: monto,

		Categoria: categoria,

		Dia: dia,

		MedioPago: medioTexto,

		Cuenta: cuenta,
		Tarjeta: tarjeta,

		AnioInicio: anio,
		MesInicio: mes,
	}

	gastosFijos = append(gastosFijos, g)

	guardarGastosFijos()

	fmt.Println("\nGasto fijo creado correctamente")

}

func verGastosFijos() {

	fmt.Println("\n--- GASTOS FIJOS ---")

	for i, g := range gastosFijos {

		fmt.Printf("%d - %s | %s | dia %d\n",

			i+1,
			g.Nombre,
			formatearMoneda(g.Monto),
			g.Dia)

	}

}