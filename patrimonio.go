package main

import "fmt"

func verPatrimonioTotal() {

    var totalARS float64
    var totalUSD float64

    fmt.Println("\n--- PATRIMONIO ---")

    fmt.Println("\nARS:")

    for _, c := range cuentas {

        if c.Moneda == "ARS" {

            fmt.Printf("%s ........ %s\n",

                c.Nombre,
                formatearMoneda(c.Saldo))

            totalARS += c.Saldo

        }

    }

    fmt.Println("\nTotal ARS:",
        formatearMoneda(totalARS))

    fmt.Println("\nUSD:")

    for _, c := range cuentas {

        if c.Moneda == "USD" {

            fmt.Printf("%s ........ USD %.2f\n",

                c.Nombre,
                c.Saldo)

            totalUSD += c.Saldo

        }

    }

    fmt.Printf("\nTotal USD: %.2f\n", totalUSD)

    var tipoCambio float64

    fmt.Println("\nTipo de cambio a usar (ej MEP):")

    fmt.Scanln(&tipoCambio)

    equivalenteARS := totalUSD * tipoCambio

    fmt.Println("\nEquivalente USD en ARS:",
        formatearMoneda(equivalenteARS))

    patrimonioTotal := totalARS + equivalenteARS

    fmt.Println("\nPATRIMONIO TOTAL:",
        formatearMoneda(patrimonioTotal))

}