package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Categoria struct {
	Nombre string
	Tipo   string
}

var categorias = []Categoria{

	{"Comida", "gasto"},
	{"Transporte", "gasto"},
	{"Servicios", "gasto"},
	{"Ocio", "gasto"},
	{"Salud", "gasto"},
	{"Educacion", "gasto"},
	{"Suscripciones", "gasto"},

	{"Sueldo", "ingreso"},
	{"Freelance", "ingreso"},
	{"Intereses", "ingreso"},
	{"Otros", "ingreso"},
}

func mostrarCategorias(tipo string) {

	println("\nSeleccione categoria:")

	numero := 1

	for _, c := range categorias {

		if c.Tipo == tipo {

			println(numero, c.Nombre)

			numero++

		}

	}
}

func elegirCategoria(tipo string) string {

	mostrarCategorias(tipo)

	var opcion int

	println("\nNumero:")
	fmt.Scanln(&opcion)

	numero := 1

	for _, c := range categorias {

		if c.Tipo == tipo {

			if numero == opcion {

				return c.Nombre

			}

			numero++

		}

	}

	println("Categoria invalida")

	return elegirCategoria(tipo)
}
func formatearMoneda(valor float64) string {

	entero := int(valor)

	texto := fmt.Sprintf("%d", entero)

	var resultado string
	contador := 0

	for i := len(texto) - 1; i >= 0; i-- {

		resultado = string(texto[i]) + resultado

		contador++

		if contador%3 == 0 && i != 0 {

			resultado = "." + resultado

		}

	}

	return "$ " + resultado

}

func guardarCategorias() {

	file, _ := os.Create("categorias.json")
	defer file.Close()

	json.NewEncoder(file).Encode(categorias)

}

func cargarCategorias() {

	file, err := os.Open("categorias.json")

	if err != nil {
		return
	}

	defer file.Close()

	json.NewDecoder(file).Decode(&categorias)

}
func verCategorias() {

	fmt.Println("\n--- CATEGORIAS ---")

	for i, c := range categorias {

		fmt.Printf("%d - %s (%s)\n",
			i+1,
			c.Nombre,
			c.Tipo)

	}

}
func agregarCategoria() {

	var nombre string
	var tipo int

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nNombre de la categoria:")
	nombre, _ = reader.ReadString('\n')

	nombre = strings.TrimSpace(nombre)

	fmt.Println("\nTipo:")
	fmt.Println("1 gasto")
	fmt.Println("2 ingreso")

	fmt.Scanln(&tipo)

	var tipoTexto string

	if tipo == 1 {

		tipoTexto = "gasto"

	} else {

		tipoTexto = "ingreso"

	}

	nueva := Categoria{
		Nombre: nombre,
		Tipo:   tipoTexto,
	}

	categorias = append(categorias, nueva)

	guardarCategorias()

	fmt.Println("Categoria agregada correctamente")

}
func eliminarCategoria() {

	verCategorias()

	var opcion int

	fmt.Println("\nNumero a eliminar:")
	fmt.Scanln(&opcion)

	if opcion <= 0 || opcion > len(categorias) {

		fmt.Println("Numero invalido")
		return

	}

	categorias = append(
		categorias[:opcion-1],
		categorias[opcion:]...)

	guardarCategorias()

	fmt.Println("Categoria eliminada")

}
