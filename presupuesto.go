package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Presupuesto struct {
	Categoria string
	Mes       int
	Monto     float64
}

var presupuestos []Presupuesto

func guardarPresupuestos() {
	file, _ := os.Create("presupuestos.json")
	defer file.Close()

	json.NewEncoder(file).Encode(presupuestos)
}

func cargarPresupuestos() {
	file, err := os.Open("presupuestos.json")

	if err != nil {
		return
	}

	defer file.Close()

	json.NewDecoder(file).Decode(&presupuestos)
}
func cargarPresupuesto() {

	var categoria string
	var mes int
	var monto float64

	println("\n--- CARGAR PRESUPUESTO ---")

	println("Categoria:")
	fmt.Scanln(&categoria)

	println("Mes (1-12):")
	fmt.Scanln(&mes)

	println("Monto presupuestado:")
	fmt.Scanln(&monto)

	p := Presupuesto{
		Categoria: categoria,
		Mes:       mes,
		Monto:     monto,
	}

	presupuestos = append(presupuestos, p)

	guardarPresupuestos()
}

func verPresupuestoPorMes() {

	var mes int

	println("\nMes a consultar:")
	fmt.Scanln(&mes)

	println("\n--- PRESUPUESTO DEL MES ---")

	for _, p := range presupuestos {

		if p.Mes == mes {

			println(p.Categoria, ":", p.Monto)

		}

	}
}
func verAnalisisPresupuesto() {

	var mes int

	fmt.Println("\nMes a analizar:")
	fmt.Scanln(&mes)

	fmt.Println("\n--- ANALISIS PRESUPUESTARIO ---")
	fmt.Println("Categoria | Presupuesto | Real | Desvio | %")

	var totalPpto float64
	var totalRealMes float64

	for _, p := range presupuestos {

		if p.Mes != mes {
			continue
		}

		var totalReal float64

		for _, m := range movimientos {

			if m.Tipo != "gasto" {
				continue
			}

			valorCuota := m.Monto / float64(m.Cuotas)

			for i := 0; i < m.Cuotas; i++ {

				mesCuota := m.Mes + i

				if mesCuota == mes && m.Categoria == p.Categoria {

					totalReal += valorCuota

				}

			}

		}

		desvio := totalReal - p.Monto

		var porcentaje float64

		if p.Monto != 0 {
			porcentaje = (desvio / p.Monto) * 100
		}

		fmt.Printf("%s | %s | %s | %s | %.1f%%\n",
			p.Categoria,
			formatearMoneda(p.Monto),
			formatearMoneda(totalReal),
			formatearMoneda(desvio),
			porcentaje)

		totalPpto += p.Monto
		totalRealMes += totalReal

	}

	totalDesvio := totalRealMes - totalPpto

	var porcentajeTotal float64

	if totalPpto != 0 {
		porcentajeTotal = (totalDesvio / totalPpto) * 100
	}

	fmt.Println("--------------------------------------------")

	fmt.Printf("TOTAL | %s | %s | %s | %.1f%%\n",
		formatearMoneda(totalPpto),
		formatearMoneda(totalRealMes),
		formatearMoneda(totalDesvio),
		porcentajeTotal)
}
func borrarPresupuestosSeguro() {

	var clave string

	fmt.Println("\nClave admin para borrar presupuestos:")
	fmt.Scanln(&clave)

	if clave != "1234" {

		fmt.Println("Clave incorrecta")
		return

	}

	presupuestos = []Presupuesto{}

	guardarPresupuestos()

	fmt.Println("Presupuestos eliminados correctamente")

}
