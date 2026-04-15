package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var movimientos []Movimiento

func guardarDatos() {
	file, _ := os.Create("data.json")
	defer file.Close()

	json.NewEncoder(file).Encode(movimientos)
}
func cargarDatos() {
	file, err := os.Open("data.json")
	if err != nil {
		return // si no existe, no pasa nada
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&movimientos)
}

func borrarMovimientosSeguro() {

	var clave string

	println("\nIngrese clave de administrador:")
	fmt.Scanln(&clave)

	if clave != "1234" {

		println("Clave incorrecta")
		return

	}

	movimientos = []Movimiento{}

	guardarDatos()

	println("Todos los movimientos fueron eliminados")

}
