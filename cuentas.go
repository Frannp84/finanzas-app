package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func guardarCuentas() {

	file, _ := os.Create("cuentas.json")
	defer file.Close()

	json.NewEncoder(file).Encode(cuentas)

}

func cargarCuentas() {

	file, err := os.Open("cuentas.json")

	if err != nil {
		return
	}

	defer file.Close()

	json.NewDecoder(file).Decode(&cuentas)

}

type Cuenta struct {

    Nombre string
    Tipo   string

    Moneda string

    Saldo float64
}

var cuentas []Cuenta

func agregarCuenta() {

    var nombre string
    var tipo int
    var moneda int

    reader := bufio.NewReader(os.Stdin)

    fmt.Println("\nNombre de la cuenta:")
    nombre, _ = reader.ReadString('\n')

    nombre = strings.TrimSpace(nombre)

    fmt.Println("\nTipo:")
    fmt.Println("1 Efectivo")
    fmt.Println("2 Banco")
    fmt.Println("3 Billetera")

    fmt.Scanln(&tipo)

    fmt.Println("\nMoneda:")
    fmt.Println("1 ARS")
    fmt.Println("2 USD")

    fmt.Scanln(&moneda)

    var tipoTexto string

    switch tipo {

    case 1:
        tipoTexto = "Efectivo"

    case 2:
        tipoTexto = "Banco"

    case 3:
        tipoTexto = "Billetera"

    }

    var monedaTexto string

    if moneda == 2 {

        monedaTexto = "USD"

    } else {

        monedaTexto = "ARS"

    }

    nueva := Cuenta{

        Nombre: nombre,
        Tipo:   tipoTexto,

        Moneda: monedaTexto,

        Saldo:  0,
    }

    cuentas = append(cuentas, nueva)

    guardarCuentas()

    fmt.Println("Cuenta creada correctamente")

}

func verCuentas() {

    fmt.Println("\n--- CUENTAS EN PESOS ---")

    for i, c := range cuentas {

        if c.Moneda == "ARS" {

            fmt.Printf("%d - %s (%s) | saldo %s\n",

                i+1,
                c.Nombre,
                c.Tipo,
                formatearMoneda(c.Saldo))
        }
    }

    fmt.Println("\n--- CUENTAS EN USD ---")

    for i, c := range cuentas {

        if c.Moneda == "USD" {

            fmt.Printf("%d - %s (%s) | saldo USD %.2f\n",

                i+1,
                c.Nombre,
                c.Tipo,
                c.Saldo)
        }
    }

}

func elegirCuenta() int {

	fmt.Println("\nSeleccione cuenta:")

	for i, c := range cuentas {

		fmt.Printf("%d %s (%s)\n",
			i+1,
			c.Nombre,
			c.Tipo)

	}

	var opcion int

	fmt.Print("\nNumero: ")
	fmt.Scanln(&opcion)

	if opcion <= 0 || opcion > len(cuentas) {

		fmt.Println("Cuenta invalida")
		return elegirCuenta()

	}

	return opcion - 1

}
