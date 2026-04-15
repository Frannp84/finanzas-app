package main

import "fmt"

func verProximosVencimientos() {

	var mes int
	var anio int

	fmt.Print("\nAño: ")
	fmt.Scanln(&anio)

	fmt.Print("Mes (1-12): ")
	fmt.Scanln(&mes)

	fmt.Println("\n--- PROXIMOS VENCIMIENTOS ---")

	// gastos fijos

	for _, g := range gastosFijos {

		if anio < g.AnioInicio {
			continue
		}

		if anio == g.AnioInicio && mes < g.MesInicio {
			continue
		}

		fmt.Printf("\n%d %s %d\n",

			g.Dia,
			nombreMes(mes),
			anio)

		fmt.Printf("%s %s\n",

			g.Nombre,
			formatearMoneda(g.Monto))

	}

	// vencimientos tarjeta

	for _, t := range tarjetas {

		fmt.Printf("\n%d %s %d\n",

			t.Vencimiento,
			nombreMes(mes),
			anio)

		fmt.Printf("Vencimiento %s\n",

			t.Nombre)

	}

}