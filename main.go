package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

func conectarDB() {

	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		log.Fatal("❌ DATABASE_URL vacío")
	}

	var err error

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error abriendo DB:", err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal("Error conectando a DB:", err)
	}

	fmt.Println("✅ DB conectada")
}

func crearTablaUsuarios() {

	query := `
	CREATE TABLE IF NOT EXISTS usuarios (
		id SERIAL PRIMARY KEY,
		usuario TEXT UNIQUE,
		password TEXT
	);
	`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatal("Error creando tabla usuarios:", err)
	}

	fmt.Println("✅ Tabla usuarios lista")
}

func crearUsuario(user string, pass string) {

	_, err := db.Exec(
		"INSERT INTO usuarios (usuario, password) VALUES ($1, $2)",
		user, pass,
	)

	if err != nil {
		fmt.Println("❌ Error creando usuario:", err)
		return
	}

	fmt.Println("✅ Usuario creado correctamente")
}

var tarjeta = Tarjeta{

	Nombre: "Visa",

	LimiteTotal:      0,
	LimiteDisponible: 0,

	Cierre:      20,
	Vencimiento: 30,

	CuentaAsociada: 0,
}

func main() {
	conectarDB()
	crearTablaUsuarios()
	crearUsuario("fran", "1234")
	cargarDatos()
	cargarPresupuestos()
	cargarCategorias()
	cargarCuentas()
	cargarTarjetas()
	cargarUsuarios()
	cargarTransferencias()
	cargarGastosFijos()
	iniciarFrontend() //

	var opcion int

	for {
		fmt.Println("\n--- APP FINANZAS ---")
		fmt.Println("1. Agregar ingreso")
		fmt.Println("2. Agregar gasto")
		fmt.Println("3. Ver balance")
		fmt.Println("4. Ver movimientos")
		fmt.Println("5. Ver gastos por categoria")
		fmt.Println("6. Ver movimientos por mes")
		fmt.Println("7. Ver resumen tarjeta")
		fmt.Println("8. Cargar presupuesto")
		fmt.Println("9. Ver presupuesto del mes")
		fmt.Println("10. Analisis presupuesto")
		fmt.Println("11. Resumen mensual")
		fmt.Println("12. Transferir entre cuentas")
		fmt.Println("13. Ver Transferencias")
		fmt.Println("14. Ver Patrimonio Total")
		fmt.Println("15. Ver Cuotas Futuras")
		fmt.Println("16. Ver deuda total tarjeta")
		fmt.Println("17. Dashboard mes actual")
		fmt.Println("18. Proximos vencimientos")
		fmt.Println("19. Ver dashboard otro periodo")
		fmt.Println("99. Configuracion")
		fmt.Println("0. Salir")
		fmt.Print("Elegí una opción: ")
		fmt.Scanln(&opcion)

		switch opcion {
		case 1:
			agregarMovimiento("ingreso")
		case 2:
			agregarMovimiento("gasto")
		case 3:
			verBalance()
		case 4:
			verMovimientos()
		case 5:
			verGastosPorCategoria()
		case 6:
			verMovimientosPorMes()
		case 7:
			verResumenTarjeta()
		case 8:
			cargarPresupuesto()
		case 9:
			verPresupuestoPorMes()
		case 10:
			verAnalisisPresupuesto()
		case 11:
			verResumenMensual()
		case 12:
			transferirEntreCuentas()
		case 13:
			verTransferencias()
		case 14:
			verPatrimonioTotal()
		case 15:
			verCuotasFuturas()
		case 16:
			verDeudaTarjeta()
		case 17:
			verDashboard()
		case 18:
			verProximosVencimientos()
		case 99:
			menuConfiguracion()
		case 0:
			fmt.Println("Hasta luego!")
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}

func agregarMovimiento(tipo string) {

	fmt.Println("\nMedio de pago:")
	fmt.Println("1 Efectivo")
	fmt.Println("2 Cuenta Bancaria")
	fmt.Println("3 Tarjeta Credito")

	var medio int
	fmt.Scanln(&medio)

	var indiceCuenta int
	var indiceTarjeta int = -1

	var medioPagoTexto string

	switch medio {

	case 1:
		medioPagoTexto = "Efectivo"
		indiceCuenta = elegirCuenta()

	case 2:
		medioPagoTexto = "Cuenta"
		indiceCuenta = elegirCuenta()

	case 3:
		medioPagoTexto = "Tarjeta"
		indiceTarjeta = elegirTarjeta()

	}

	var monto float64
	var detalle string
	var Categoria string
	var Anio int
	var Mes int
	var Dia int
	var Cuotas int

	fmt.Print("Monto: ")
	fmt.Scanln(&monto)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Detalle: ")
	detalle, _ = reader.ReadString('\n')
	detalle = strings.TrimSpace(detalle)

	Categoria = elegirCategoria(tipo)

	fmt.Print("Año: ")
	fmt.Scanln(&Anio)

	fmt.Print("Mes (1-12): ")
	fmt.Scanln(&Mes)

	fmt.Print("Dia del gasto (1-31): ")
	fmt.Scanln(&Dia)

	if tipo == "gasto" {

		if medioPagoTexto == "Tarjeta" {

			fmt.Print("Cantidad de cuotas: ")
			fmt.Scanln(&Cuotas)

			if Cuotas <= 0 {
				Cuotas = 1
			}

		} else {

			Cuotas = 1

		}

	} else {

		Cuotas = 1

	}

	mov := Movimiento{

		Tipo:  tipo,
		Monto: monto,

		Detalle:   detalle,
		Categoria: Categoria,

		Anio: Anio,
		Mes:  Mes,
		Dia:  Dia,

		Cuotas: Cuotas,

		MedioPago: medioPagoTexto,

		Cuenta:  indiceCuenta,
		Tarjeta: indiceTarjeta,
	}

	movimientos = append(movimientos, mov)

	valorCuota := mov.Monto / float64(mov.Cuotas)

	// impacto inmediato solo si NO es tarjeta
	if medioPagoTexto != "Tarjeta" {

		if tipo == "ingreso" {

			cuentas[indiceCuenta].Saldo += valorCuota

		} else {

			cuentas[indiceCuenta].Saldo -= valorCuota

		}

	}

	// consumir limite de tarjeta
	if medioPagoTexto == "Tarjeta" && tipo == "gasto" {

		tarjetas[indiceTarjeta].LimiteDisponible -= mov.Monto

	}

	guardarCuentas()
	guardarTarjetas()
	guardarDatos()

}

func verBalance() {
	var balance float64
	var mesActual int

	fmt.Print("Ingrese mes a analizar (1-12): ")
	fmt.Scanln(&mesActual)

	for _, m := range movimientos {
		valorCuota := m.Monto / float64(m.Cuotas)

		for i := 0; i < m.Cuotas; i++ {
			mesCuota := m.Mes + i

			if mesCuota == mesActual {
				if m.Tipo == "ingreso" {
					balance += valorCuota
				} else {
					balance -= valorCuota
				}
			}
		}
	}

	fmt.Println("Balance del mes:", formatearMoneda(balance))
}

func verMovimientos() {

	fmt.Println("\n--- TODOS LOS MOVIMIENTOS ---")

	for i, m := range movimientos {

		fmt.Printf("%d - %s | %s | %s | %s | %d %s %d | Cuotas: %d\n",

			i+1,
			m.Tipo,
			formatearMoneda(m.Monto),
			m.Detalle,
			m.Categoria,

			m.Dia,
			nombreMes(m.Mes),
			m.Anio,

			m.Cuotas)
	}

}
func verMovimientosPorMes() {

	var mesActual int
	var anioActual int

	fmt.Print("Ingrese año: ")
	fmt.Scanln(&anioActual)

	fmt.Print("Ingrese mes (1-12): ")
	fmt.Scanln(&mesActual)

	fmt.Println("\n--- MOVIMIENTOS DEL MES ---")

	for i, m := range movimientos {

		valorCuota := m.Monto / float64(m.Cuotas)

		for j := 0; j < m.Cuotas; j++ {

			mesCuota := m.Mes + j
			anioCuota := m.Anio

			for mesCuota > 12 {

				mesCuota -= 12
				anioCuota++

			}

			if mesCuota == mesActual && anioCuota == anioActual {

				if m.Tipo == "ingreso" {

					fmt.Printf("%d - 💰 ingreso | %s | %s | %s\n",
						i+1,
						formatearMoneda(valorCuota),
						m.Detalle,
						m.Categoria)

				} else {

					fmt.Printf("%d - 💸 gasto | %s | %s | %s\n",
						i+1,
						formatearMoneda(valorCuota),
						m.Detalle,
						m.Categoria)

				}

			}

		}

	}

}
func verGastosPorCategoria() {
	gastos := make(map[string]float64)

	for _, m := range movimientos {
		if m.Tipo == "gasto" {
			gastos[m.Categoria] += m.Monto
		}
	}

	fmt.Println("\n--- GASTOS POR CATEGORIA ---")
	for categoria, total := range gastos {
		fmt.Printf("%s: %.2f\n", categoria, total)
	}
}

func verResumenTarjeta() {

	if len(tarjetas) == 0 {

		fmt.Println("No hay tarjetas cargadas")
		return

	}

	indiceTarjeta := elegirTarjeta()

	t := tarjetas[indiceTarjeta]

	var mes int

	fmt.Print("\nIngrese mes del resumen: ")
	fmt.Scanln(&mes)

	fmt.Println("\n--- RESUMEN TARJETA ---")

	fmt.Println("\nTarjeta:", t.Nombre)

	var total float64

	// MOVIMIENTOS NORMALES

	for _, m := range movimientos {

		if m.Tipo != "gasto" {
			continue
		}

		if m.Tarjeta != indiceTarjeta {
			continue
		}

		valorCuota := m.Monto / float64(m.Cuotas)

		mesResumen := calcularMesResumen(m.Mes, m.Dia, t.Cierre)

		for i := 0; i < m.Cuotas; i++ {

			if mesResumen+i == mes {

				total += valorCuota

				fmt.Printf("💸 %s | %s\n",
					formatearMoneda(valorCuota),
					m.Detalle)

			}

		}

	}

	// GASTOS FIJOS

	for _, g := range gastosFijos {

		if g.MedioPago != "Tarjeta" {
			continue
		}

		if g.Tarjeta != indiceTarjeta {
			continue
		}

		if mes < g.MesInicio {
			continue
		}

		mesResumen := calcularMesResumen(mes, g.Dia, t.Cierre)

		if mesResumen == mes {

			total += g.Monto

			fmt.Printf("💸 %s | %s (fijo)\n",
				formatearMoneda(g.Monto),
				g.Nombre)

		}

	}

	fmt.Println("\nTotal a pagar:", formatearMoneda(total))

	fmt.Printf("Vence el día %d\n", t.Vencimiento)

}

func verResumenMensual() {

	var mes int
	var anio int

	generarMovimientosGastosFijos(mes, anio)

	fmt.Print("\nAño: ")
	fmt.Scanln(&anio)

	fmt.Print("Mes (1-12): ")
	fmt.Scanln(&mes)

	var totalIngresos float64
	var totalGastos float64

	gastosPorCategoria := make(map[string]float64)

	var totalCuotasMes float64

	for _, g := range gastosFijos {

		if anio > g.AnioInicio ||
			(anio == g.AnioInicio && mes >= g.MesInicio) {

			if g.MedioPago == "Tarjeta" {

				t := tarjetas[g.Tarjeta]

				mesResumen := calcularMesResumen(mes, g.Dia, t.Cierre)

				if mesResumen == mes {

					totalGastos += g.Monto
					gastosPorCategoria[g.Categoria] += g.Monto

				}

			} else {

				totalGastos += g.Monto
				gastosPorCategoria[g.Categoria] += g.Monto

			}

		}
	}

	for _, m := range movimientos {

		valorCuota := m.Monto / float64(m.Cuotas)

		for i := 0; i < m.Cuotas; i++ {

			mesCuota := m.Mes + i
			anioCuota := m.Anio

			for mesCuota > 12 {

				mesCuota -= 12
				anioCuota++

			}

			if mesCuota == mes && anioCuota == anio {

				if m.Tipo == "ingreso" {

					totalIngresos += valorCuota

				} else {

					totalGastos += valorCuota

					gastosPorCategoria[m.Categoria] += valorCuota

					if m.Cuotas > 1 {

						totalCuotasMes += valorCuota

					}

				}

			}

		}

	}

	resultado := totalIngresos - totalGastos

	var porcentajeAhorro float64

	if totalIngresos > 0 {

		porcentajeAhorro = (resultado / totalIngresos) * 100

	}

	// detectar categoria con mayor gasto

	var categoriaMayor string
	var mayorGasto float64

	for cat, total := range gastosPorCategoria {

		if total > mayorGasto {

			mayorGasto = total
			categoriaMayor = cat

		}

	}

	fmt.Println("\n--- RESUMEN MENSUAL ---")

	fmt.Printf("\nPeriodo: %s %d\n",

		nombreMes(mes),
		anio)

	fmt.Println("\nIngresos:",
		formatearMoneda(totalIngresos))

	fmt.Println("Gastos:",
		formatearMoneda(totalGastos))

	fmt.Println("\nResultado:",
		formatearMoneda(resultado))

	fmt.Printf("\nAhorro: %.1f%%\n",
		porcentajeAhorro)

	fmt.Println("\nCategoria con mayor gasto:",
		categoriaMayor,
		formatearMoneda(mayorGasto))

	fmt.Println("\nCuotas del mes:",
		formatearMoneda(totalCuotasMes))

}

func verTransferencias() {

	fmt.Println("\n--- HISTORIAL DE TRANSFERENCIAS ---")

	for i, t := range transferencias {

		fmt.Printf("%d | %s → %s | %s | %d/%d\n",

			i+1,

			cuentas[t.Origen].Nombre,
			cuentas[t.Destino].Nombre,

			formatearMoneda(t.Monto),

			t.Dia,
			t.Mes,
		)
	}

}
