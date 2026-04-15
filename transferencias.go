package main

import "fmt"

func transferirEntreCuentas() {

    if len(cuentas) < 2 {

        fmt.Println("\nNecesitas al menos 2 cuentas")
        return

    }

    fmt.Println("\n--- TRANSFERENCIA ---")

    fmt.Println("\nCuenta origen:")
    origen := elegirCuenta()

    fmt.Println("\nCuenta destino:")
    destino := elegirCuenta()

    if origen == destino {

        fmt.Println("No se puede transferir a la misma cuenta")
        return

    }
	if cuentas[origen].Moneda != cuentas[destino].Moneda {

    fmt.Println("\nNo se puede transferir entre monedas distintas")
    fmt.Println("Proximamente conversion automatica")

    return
}

    var monto float64
    var mes int
    var dia int

    fmt.Println("\nMonto:")
    fmt.Scanln(&monto)

    fmt.Println("\nMes (1-12):")
    fmt.Scanln(&mes)

    fmt.Println("\nDia:")
    fmt.Scanln(&dia)

    if cuentas[origen].Saldo < monto {

        fmt.Println("Saldo insuficiente")
        return

    }

    cuentas[origen].Saldo -= monto
    cuentas[destino].Saldo += monto

    t := Transferencia{

        Origen: origen,
        Destino: destino,

        Monto: monto,

        Mes: mes,
        Dia: dia,
    }

    transferencias = append(transferencias, t)

    guardarCuentas()
    guardarTransferencias()

    fmt.Println("\nTransferencia registrada correctamente")
}
