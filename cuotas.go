package main

import "fmt"

func verCuotasFuturas() {

    cuotasPorPeriodo := make(map[string]float64)

    for _, m := range movimientos {

        if m.Tipo != "gasto" {
            continue
        }

        if m.Cuotas <= 1 {
            continue
        }

        valorCuota := m.Monto / float64(m.Cuotas)

        for i := 1; i < m.Cuotas; i++ {

            mesCuota := m.Mes + i
            anioCuota := m.Anio

            for mesCuota > 12 {

                mesCuota -= 12
                anioCuota++

            }

            clave := fmt.Sprintf("%s %d",

                nombreMes(mesCuota),
                anioCuota)

            cuotasPorPeriodo[clave] += valorCuota

        }

    }

    fmt.Println("\n--- CUOTAS FUTURAS ---")

    var totalFuturo float64

    for periodo, total := range cuotasPorPeriodo {

        fmt.Printf("%s: %s\n",

            periodo,
            formatearMoneda(total))

        totalFuturo += total

    }

    fmt.Println("\nTOTAL COMPROMETIDO:",
        formatearMoneda(totalFuturo))

}