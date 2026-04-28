package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func requiereLogin(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if usuarioActual == "" {

			http.Redirect(w, r, "/login", http.StatusSeeOther)

			return

		}

		handler(w, r)

	}

}

func iniciarFrontend() {

	// DASHBOARD
	cargarUsuarios()
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/registro", registroHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if usuarioActual == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		now := time.Now()

		mes := int(now.Month())
		anio := now.Year()

		mesURL := r.URL.Query().Get("mes")
		anioURL := r.URL.Query().Get("anio")

		if mesURL != "" {

			m, err := strconv.Atoi(mesURL)

			if err == nil && m >= 1 && m <= 12 {

				mes = m

			}

		}

		if anioURL != "" {

			a, err := strconv.Atoi(anioURL)

			if err == nil {

				anio = a

			}

		}

		generarMovimientosGastosFijos(mes, anio)

		var totalARS float64
		var totalUSD float64

		for _, c := range cuentas {

			if c.Moneda == "ARS" {

				totalARS += c.Saldo

			} else {

				totalUSD += c.Saldo

			}

		}

		var totalIngresos float64
		var totalGastos float64

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

					}

				}

			}

		}

		resultado := totalIngresos - totalGastos

		var porcentajeAhorro float64

		if totalIngresos > 0 {

			porcentajeAhorro = (resultado / totalIngresos) * 100

		}

		var cuotasFuturas float64

		for _, m := range movimientos {

			if m.Tipo != "gasto" {

				continue

			}

			if m.Cuotas <= 1 {

				continue

			}

			valorCuota := m.Monto / float64(m.Cuotas)

			for i := 1; i < m.Cuotas; i++ {

				cuotasFuturas += valorCuota

			}

		}

		var deudaTarjetas float64

		for _, t := range tarjetas {

			deudaTarjetas += (t.LimiteTotal - t.LimiteDisponible)

		}

		fmt.Fprintf(w, `
<html>

<head>

<title>Dashboard</title>

<style>

body {

	font-family: Arial;
	background:#0f172a;
	color:white;
	padding:40px;

}

.card {

	background:#1e293b;
	padding:20px;
	border-radius:10px;
	width:260px;
	margin:10px;
	display:inline-block;
	vertical-align:top;

}

button {

	padding:10px;
	margin:5px;
	border:none;
	border-radius:6px;
	cursor:pointer;

}

select, input {

	padding:6px;
	margin:4px;

}

</style>

</head>

<body>

<h1>📊 Dashboard financiero</h1>

<br>

<button onclick="location.href='/logout'">

Cerrar sesión

</button>

<form>

Mes:

<select name="mes">

<option value="1">Enero</option>
<option value="2">Febrero</option>
<option value="3">Marzo</option>
<option value="4">Abril</option>
<option value="5">Mayo</option>
<option value="6">Junio</option>
<option value="7">Julio</option>
<option value="8">Agosto</option>
<option value="9">Septiembre</option>
<option value="10">Octubre</option>
<option value="11">Noviembre</option>
<option value="12">Diciembre</option>

</select>

Año:

<input name="anio" value="%d" style="width:90px">

<button type="submit">

Ver

</button>

</form>

<br>

<div class="card">

<h3>Patrimonio ARS</h3>

<h2>%s</h2>

</div>

<div class="card">

<h3>Ingresos del mes</h3>

<h2>%s</h2>

</div>

<div class="card">

<h3>Gastos del mes</h3>

<h2>%s</h2>

</div>

<div class="card">

<h3>Resultado mensual</h3>

<h2>%s</h2>

</div>

<div class="card">

<h3>Ahorro</h3>

<h2>%.1f%%</h2>

</div>

<div class="card">

<h3>Cuotas futuras</h3>

<h2>%s</h2>

</div>

<div class="card">

<h3>Deuda tarjetas</h3>

<h2>%s</h2>

</div>

<br><br>

<button onclick="location.href='/movimientos'">Movimientos</button>

<button onclick="location.href='/agregar'">Agregar movimiento</button>

<button onclick="location.href='/vencimientos'">Vencimientos</button>

<button onclick="location.href='/categorias'">Gastos por categoria</button>

<button onclick="location.href='/tarjetas'">Tarjetas</button>

<button onclick="location.href='/transferencias'">Transferencias</button>

<button onclick="location.href='/gastosfijos'">Gastos fijos</button>

<button onclick="location.href='/presupuesto'">Analisis presupuesto</button>

<button onclick="location.href='/crear_presupuesto'">Crear presupuesto</button>

</body>

</html>
`,
			anio,
			formatearMoneda(totalARS),
			formatearMoneda(totalIngresos),
			formatearMoneda(totalGastos),
			formatearMoneda(resultado),
			porcentajeAhorro,
			formatearMoneda(cuotasFuturas),
			formatearMoneda(deudaTarjetas),
		)

	})

	// MOVIMIENTOS
	http.HandleFunc("/movimientos", requiereLogin(func(w http.ResponseWriter, r *http.Request) {

		if usuarioActual == "" {

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return

		}

		fmt.Fprintf(w, `

		<html>

		<head>

		<style>

		body {

			font-family: Arial;

			background:#0f172a;

			color:white;

			padding:40px;

		}

		table {

			border-collapse: collapse;

			width:90%%;

		}

		td, th {

			padding:8px;

			border-bottom:1px solid gray;

		}

		button {

			padding:10px;

			margin-top:20px;

		}

		</style>

		</head>

		<body>

		<h1>Movimientos</h1>

		<table>

		<tr>

		<th>Tipo</th>
		<th>Monto</th>
		<th>Categoria</th>
		<th>Fecha</th>
		<th>Cuotas</th>

		</tr>

		`)

		for _, m := range movimientos {

			fmt.Fprintf(w,

				`<tr>

				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%d/%d/%d</td>
				<td>%d</td>

				</tr>`,

				m.Tipo,
				formatearMoneda(m.Monto),
				m.Categoria,
				m.Dia,
				m.Mes,
				m.Anio,
				m.Cuotas,
			)

		}

		fmt.Fprintf(w, `

		</table>

		<br>

		<button onclick="location.href='/'">

		Volver

		</button>

		</body>

		</html>

		`)

	}))

	// FORMULARIO COMPLETO
	http.HandleFunc("/agregar", requiereLogin(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {

			tipo := r.FormValue("tipo")

			categoria := r.FormValue("categoria")

			detalle := r.FormValue("detalle")

			montoTexto := r.FormValue("monto")

			anioTexto := r.FormValue("anio")

			mesTexto := r.FormValue("mes")

			diaTexto := r.FormValue("dia")

			cuentaTexto := r.FormValue("cuenta")

			tarjetaTexto := r.FormValue("tarjeta")

			cuotasTexto := r.FormValue("cuotas")

			monto, _ := strconv.ParseFloat(montoTexto, 64)

			anio, _ := strconv.Atoi(anioTexto)

			mes, _ := strconv.Atoi(mesTexto)

			dia, _ := strconv.Atoi(diaTexto)

			cuenta, _ := strconv.Atoi(cuentaTexto)

			tarjeta, _ := strconv.Atoi(tarjetaTexto)

			cuotas, _ := strconv.Atoi(cuotasTexto)

			if cuotas <= 0 {

				cuotas = 1

			}

			mov := Movimiento{

				Tipo: tipo,

				Monto: monto,

				Detalle: detalle,

				Categoria: categoria,

				Anio: anio,

				Mes: mes,

				Dia: dia,

				Cuotas: cuotas,

				Cuenta: cuenta,

				Tarjeta: tarjeta,
			}

			movimientos = append(movimientos, mov)

			valorCuota := monto / float64(cuotas)

			if tarjeta == -1 {

				if tipo == "ingreso" {

					cuentas[cuenta].Saldo += valorCuota

				} else {

					cuentas[cuenta].Saldo -= valorCuota

				}

			}

			if tarjeta >= 0 {

				tarjetas[tarjeta].LimiteDisponible -= monto

			}

			guardarDatos()

			guardarCuentas()

			guardarTarjetas()

			http.Redirect(w, r, "/", http.StatusSeeOther)

			return

		}

		// cargar opciones dinamicas

		tipoSeleccionado := r.URL.Query().Get("tipo")

		if tipoSeleccionado == "" {

			tipoSeleccionado = "gasto"

		}

		var opcionesCategorias string

		for _, c := range categorias {

			if c.Tipo == tipoSeleccionado {

				opcionesCategorias += fmt.Sprintf(

					"<option value='%s'>%s</option>",

					c.Nombre,

					c.Nombre,
				)

			}

		}

		var opcionesCuentas string

		for i, c := range cuentas {

			opcionesCuentas += fmt.Sprintf(

				"<option value='%d'>%s</option>",

				i,

				c.Nombre,
			)

		}

		var opcionesTarjetas string

		opcionesTarjetas = "<option value='-1'>No usar tarjeta</option>"

		for i, t := range tarjetas {

			opcionesTarjetas += fmt.Sprintf(

				"<option value='%d'>%s</option>",

				i,

				t.Nombre,
			)

		}

		fmt.Fprintf(w, `

		<html>

		<head>

		<style>

		body {

			font-family: Arial;

			background:#0f172a;

			color:white;

			padding:40px;

		}

		input, select {

			padding:8px;

			margin:5px;

		}

		button {

			padding:10px;

			margin-top:10px;

		}

		</style>

		</head>

		<body>

		<h1>Agregar movimiento</h1>

		<form method="POST">

		Tipo:

		<select name="tipo" onchange="window.location='?tipo='+this.value">

		<option value="gasto">Gasto</option>

		<option value="ingreso">Ingreso</option>

		</select>

		<br>

		Categoria:

		<select name="categoria">

		%s

		</select>

		<br>

		Cuenta:

		<select name="cuenta">

		%s

		</select>

		<br>

		Tarjeta:

		<select name="tarjeta">

		%s

		</select>

		<br>

		Monto:

		<input name="monto">

		<br>

		Detalle:

		<input name="detalle">

		<br>

		Año:

		<input name="anio" value="2026">

		<br>

		Mes:

		<input name="mes" value="3">

		<br>

		Dia:

		<input name="dia" value="25">

		<br>

		Cuotas:

		<input name="cuotas" value="1">

		<br><br>

		<button type="submit">

		Guardar

		</button>

		</form>

		<br>

		<button onclick="location.href='/'">

		Volver

		</button>

		</body>

		</html>

		`,

			opcionesCategorias,

			opcionesCuentas,

			opcionesTarjetas,
		)

	}))

	// VENCIMIENTOS
	http.HandleFunc("/vencimientos", requiereLogin(func(w http.ResponseWriter, r *http.Request) {

		var mes int = 3
		var anio int = 2026

		fmt.Fprintf(w, `

	<html>

	<head>

	<style>

	body {

		font-family: Arial;

		background:#0f172a;

		color:white;

		padding:40px;

	}

	table {

		border-collapse: collapse;

		width:70%%;

	}

	td, th {

		padding:8px;

		border-bottom:1px solid gray;

	}

	button {

		padding:10px;

		margin-top:20px;

	}

	</style>

	</head>

	<body>

	<h1>📅 Próximos vencimientos</h1>

	<table>

	<tr>

	<th>Fecha</th>
	<th>Concepto</th>
	<th>Monto</th>

	</tr>

	`)

		// GASTOS FIJOS

		for _, g := range gastosFijos {

			if anio < g.AnioInicio {

				continue

			}

			if anio == g.AnioInicio && mes < g.MesInicio {

				continue

			}

			fmt.Fprintf(w,

				`<tr>

			<td>%d/%d/%d</td>

			<td>%s</td>

			<td>%s</td>

			</tr>`,

				g.Dia,

				mes,

				anio,

				g.Nombre,

				formatearMoneda(g.Monto),
			)

		}

		// TARJETAS

		for _, t := range tarjetas {

			fmt.Fprintf(w,

				`<tr>

			<td>%d/%d/%d</td>

			<td>Vencimiento %s</td>

			<td>-</td>

			</tr>`,

				t.Vencimiento,

				mes,

				anio,

				t.Nombre,
			)

		}

		fmt.Fprintf(w, `

	</table>

	<br>

	<button onclick="location.href='/'">

	Volver

	</button>

	</body>

	</html>

	`)

	}))

	// GASTOS POR CATEGORIA
	http.HandleFunc("/categorias", requiereLogin(func(w http.ResponseWriter, r *http.Request) {

		gastos := make(map[string]float64)

		for _, m := range movimientos {

			if m.Tipo == "gasto" {

				gastos[m.Categoria] += m.Monto

			}

		}

		fmt.Fprintf(w, `

	<html>

	<head>

	<style>

	body {

		font-family: Arial;

		background:#0f172a;

		color:white;

		padding:40px;

	}

	table {

		border-collapse: collapse;

		width:60%%;

	}

	td {

		padding:8px;

		border-bottom:1px solid gray;

	}

	.bar {

		background:#38bdf8;

		height:10px;

		border-radius:5px;

	}

	button {

		padding:10px;

		margin-top:20px;

	}

	</style>

	</head>

	<body>

	<h1>📊 Gastos por categoria</h1>

	<table>

	`)

		var max float64

		for _, total := range gastos {

			if total > max {

				max = total

			}

		}

		for cat, total := range gastos {

			porcentaje := (total / max) * 100

			fmt.Fprintf(w,

				`

			<tr>

			<td style="width:150px">%s</td>

			<td style="width:120px">%s</td>

			<td>

			<div class="bar" style="width:%.0f%%"></div>

			</td>

			</tr>

			`,

				cat,

				formatearMoneda(total),

				porcentaje,
			)

		}

		fmt.Fprintf(w, `

	</table>

	<br>

	<button onclick="location.href='/'">

	Volver

	</button>

	</body>

	</html>

	`)

	}))

	// TARJETAS
	http.HandleFunc("/tarjetas", requiereLogin(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, `

	<html>

	<head>

	<style>

	body {

		font-family: Arial;

		background:#0f172a;

		color:white;

		padding:40px;

	}

	.card {

		background:#1e293b;

		padding:20px;

		border-radius:10px;

		margin:10px;

		width:280px;

		display:inline-block;

	}

	button {

		padding:10px;

		margin-top:20px;

	}

	</style>

	</head>

	<body>

	<h1>💳 Tarjetas</h1>

	`)

		for _, t := range tarjetas {

			deuda := t.LimiteTotal - t.LimiteDisponible

			fmt.Fprintf(w,

				`

			<div class="card">

			<h3>%s</h3>

			Limite total:<br>
			%s

			<br><br>

			Disponible:<br>
			%s

			<br><br>

			Deuda:<br>
			%s

			<br><br>

			Cierre: %d<br>

			Vencimiento: %d

			</div>

			`,

				t.Nombre,

				formatearMoneda(t.LimiteTotal),

				formatearMoneda(t.LimiteDisponible),

				formatearMoneda(deuda),

				t.Cierre,

				t.Vencimiento,
			)

		}

		fmt.Fprintf(w, `

	<br>

	<button onclick="location.href='/'">

	Volver

	</button>

	</body>

	</html>

	`)

	}))

	// TRANSFERENCIAS
	http.HandleFunc("/transferencias", requiereLogin(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {

			origenTexto := r.FormValue("origen")

			destinoTexto := r.FormValue("destino")

			montoTexto := r.FormValue("monto")

			mesTexto := r.FormValue("mes")

			diaTexto := r.FormValue("dia")

			origen, _ := strconv.Atoi(origenTexto)

			destino, _ := strconv.Atoi(destinoTexto)

			monto, _ := strconv.ParseFloat(montoTexto, 64)

			mes, _ := strconv.Atoi(mesTexto)

			dia, _ := strconv.Atoi(diaTexto)

			if origen != destino {

				if cuentas[origen].Saldo >= monto {

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

				}

			}

			http.Redirect(w, r, "/", http.StatusSeeOther)

			return

		}

		var opciones string

		for i, c := range cuentas {

			opciones += fmt.Sprintf(

				"<option value='%d'>%s</option>",

				i,

				c.Nombre,
			)

		}

		fmt.Fprintf(w, `

	<html>

	<head>

	<style>

	body {

		font-family: Arial;

		background:#0f172a;

		color:white;

		padding:40px;

	}

	input, select {

		padding:8px;

		margin:5px;

	}

	button {

		padding:10px;

	}

	</style>

	</head>

	<body>

	<h1>💸 Transferencia</h1>

	<form method="POST">

	Cuenta origen:

	<select name="origen">

	%s

	</select>

	<br>

	Cuenta destino:

	<select name="destino">

	%s

	</select>

	<br>

	Monto:

	<input name="monto">

	<br>

	Mes:

	<input name="mes" value="3">

	<br>

	Dia:

	<input name="dia" value="25">

	<br><br>

	<button type="submit">

	Transferir

	</button>

	</form>

	<br>

	<button onclick="location.href='/'">

	Volver

	</button>

	</body>

	</html>

	`,

			opciones,

			opciones,
		)

	}))

	// GASTOS FIJOS
	http.HandleFunc("/gastosfijos", requiereLogin(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {

			nombre := r.FormValue("nombre")

			montoTexto := r.FormValue("monto")

			diaTexto := r.FormValue("dia")

			categoria := r.FormValue("categoria")

			cuentaTexto := r.FormValue("cuenta")

			tarjetaTexto := r.FormValue("tarjeta")

			mesTexto := r.FormValue("mes")

			anioTexto := r.FormValue("anio")

			monto, _ := strconv.ParseFloat(montoTexto, 64)

			dia, _ := strconv.Atoi(diaTexto)

			cuenta, _ := strconv.Atoi(cuentaTexto)

			tarjeta, _ := strconv.Atoi(tarjetaTexto)

			mes, _ := strconv.Atoi(mesTexto)

			anio, _ := strconv.Atoi(anioTexto)

			medio := "Cuenta"

			if tarjeta >= 0 {

				medio = "Tarjeta"

			}

			g := GastoFijo{

				Nombre: nombre,

				Monto: monto,

				Categoria: categoria,

				Dia: dia,

				MedioPago: medio,

				Cuenta: cuenta,

				Tarjeta: tarjeta,

				MesInicio: mes,

				AnioInicio: anio,
			}

			gastosFijos = append(gastosFijos, g)

			guardarGastosFijos()

			http.Redirect(w, r, "/", http.StatusSeeOther)

			return

		}

		var opcionesCategorias string

		for _, c := range categorias {

			if c.Tipo == "gasto" {

				opcionesCategorias += fmt.Sprintf(

					"<option value='%s'>%s</option>",

					c.Nombre,

					c.Nombre,
				)

			}

		}

		var opcionesCuentas string

		for i, c := range cuentas {

			opcionesCuentas += fmt.Sprintf(

				"<option value='%d'>%s</option>",

				i,

				c.Nombre,
			)

		}

		var opcionesTarjetas string

		opcionesTarjetas = "<option value='-1'>No usar tarjeta</option>"

		for i, t := range tarjetas {

			opcionesTarjetas += fmt.Sprintf(

				"<option value='%d'>%s</option>",

				i,

				t.Nombre,
			)

		}

		fmt.Fprintf(w, `

	<html>

	<head>

	<style>

	body {

		font-family: Arial;

		background:#0f172a;

		color:white;

		padding:40px;

	}

	input, select {

		padding:8px;

		margin:5px;

	}

	button {

		padding:10px;

	}

	</style>

	</head>

	<body>

	<h1>📅 Gasto fijo</h1>

	<form method="POST">

	Nombre:

	<input name="nombre">

	<br>

	Monto:

	<input name="monto">

	<br>

	Categoria:

	<select name="categoria">

	%s

	</select>

	<br>

	Dia del mes:

	<input name="dia">

	<br>

	Cuenta:

	<select name="cuenta">

	%s

	</select>

	<br>

	Tarjeta:

	<select name="tarjeta">

	%s

	</select>

	<br>

	Mes inicio:

	<input name="mes" value="3">

	<br>

	Año inicio:

	<input name="anio" value="2026">

	<br><br>

	<button type="submit">

	Guardar gasto fijo

	</button>

	</form>

	<br>

	<button onclick="location.href='/'">

	Volver

	</button>

	</body>

	</html>

	`,

			opcionesCategorias,

			opcionesCuentas,

			opcionesTarjetas,
		)

	}))

	// ANALISIS PRESUPUESTARIO
	http.HandleFunc("/presupuesto", requiereLogin(func(w http.ResponseWriter, r *http.Request) {

		var mes int = 3

		fmt.Fprintf(w, `

	<html>

	<head>

	<style>

	body {

		font-family: Arial;

		background:#0f172a;

		color:white;

		padding:40px;

	}

	table {

		border-collapse: collapse;

		width:80%%;

	}

	td, th {

		padding:8px;

		border-bottom:1px solid gray;

	}

	.pos {

		color:#22c55e;

	}

	.neg {

		color:#ef4444;

	}

	button {

		padding:10px;

		margin-top:20px;

	}

	</style>

	</head>

	<body>

	<h1>📊 Analisis presupuestario</h1>

	<table>

	<tr>

	<th>Categoria</th>
	<th>Presupuesto</th>
	<th>Real</th>
	<th>Diferencia</th>
	<th>%%</th>

	</tr>

	`)

		var totalPpto float64
		var totalReal float64

		for _, p := range presupuestos {

			if p.Mes != mes {

				continue

			}

			var real float64

			for _, m := range movimientos {

				if m.Tipo != "gasto" {

					continue

				}

				valorCuota := m.Monto / float64(m.Cuotas)

				for i := 0; i < m.Cuotas; i++ {

					mesCuota := m.Mes + i

					if mesCuota == mes && m.Categoria == p.Categoria {

						real += valorCuota

					}

				}

			}

			diferencia := real - p.Monto

			var porcentaje float64

			if p.Monto != 0 {

				porcentaje = (diferencia / p.Monto) * 100

			}

			clase := "pos"

			if diferencia > 0 {

				clase = "neg"

			}

			fmt.Fprintf(w,

				`

			<tr>

			<td>%s</td>

			<td>%s</td>

			<td>%s</td>

			<td class="%s">%s</td>

			<td>%.1f%%</td>

			</tr>

			`,

				p.Categoria,

				formatearMoneda(p.Monto),

				formatearMoneda(real),

				clase,

				formatearMoneda(diferencia),

				porcentaje,
			)

			totalPpto += p.Monto

			totalReal += real

		}

		totalDif := totalReal - totalPpto

		var totalPct float64

		if totalPpto != 0 {

			totalPct = (totalDif / totalPpto) * 100

		}

		fmt.Fprintf(w,

			`

		<tr>

		<th>TOTAL</th>

		<th>%s</th>

		<th>%s</th>

		<th>%s</th>

		<th>%.1f%%</th>

		</tr>

		`,

			formatearMoneda(totalPpto),

			formatearMoneda(totalReal),

			formatearMoneda(totalDif),

			totalPct,
		)

		fmt.Fprintf(w, `

	</table>

	<br>

	<button onclick="location.href='/'">

	Volver

	</button>

	</body>

	</html>

	`)

	}))

	// CREAR PRESUPUESTO
	http.HandleFunc("/crear_presupuesto", requiereLogin(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {

			categoria := r.FormValue("categoria")

			mesTexto := r.FormValue("mes")

			montoTexto := r.FormValue("monto")

			mes, _ := strconv.Atoi(mesTexto)

			monto, _ := strconv.ParseFloat(montoTexto, 64)

			p := Presupuesto{

				Categoria: categoria,

				Mes: mes,

				Monto: monto,
			}

			presupuestos = append(presupuestos, p)

			guardarPresupuestos()

			http.Redirect(w, r, "/presupuesto", http.StatusSeeOther)

			return

		}

		var opciones string

		for _, c := range categorias {

			if c.Tipo == "gasto" {

				opciones += fmt.Sprintf(

					"<option value='%s'>%s</option>",

					c.Nombre,

					c.Nombre,
				)

			}

		}

		fmt.Fprintf(w, `

	<html>

	<head>

	<style>

	body {

		font-family: Arial;

		background:#0f172a;

		color:white;

		padding:40px;

	}

	input, select {

		padding:8px;

		margin:5px;

	}

	button {

		padding:10px;

	}

	</style>

	</head>

	<body>

	<h1>📊 Crear presupuesto</h1>

	<form method="POST">

	Categoria:

	<select name="categoria">

	%s

	</select>

	<br>

	Mes:

	<input name="mes" value="3">

	<br>

	Monto:

	<input name="monto">

	<br><br>

	<button type="submit">

	Guardar presupuesto

	</button>

	</form>

	<br>

	<button onclick="location.href='/presupuesto'">

	Ver analisis

	</button>

	</body>

	</html>

	`, opciones)

	}))

	fmt.Println("Frontend corriendo en http://localhost:8080")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)

}
func mostrarLogin(w http.ResponseWriter) {

	fmt.Fprintf(w, `
	<html>

	<head>

	<title>Login</title>

	<style>

	body {

		font-family: Arial;
		background:#0f172a;
		color:white;
		padding:40px;

	}

	input {

		padding:10px;
		margin:10px;

	}

	button {

		padding:10px;

	}

	</style>

	</head>

	<body>

	<h1>Login</h1>

	<form method="POST" action="/login">

	Usuario:

	<br>

	<input name="usuario">

	<br>

	Password:

	<br>

	<input type="password" name="password">

	<br>

	<button>

	Ingresar

	</button>

	</form>

	</body>

	</html>
	`)

}

var usuarioActual string

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		mostrarLogin(w)

		return

	}

	user := r.FormValue("usuario")

	pass := r.FormValue("password")

	if validarUsuario(user, pass) {

		usuarioActual = user

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return

	}

	fmt.Fprintf(w, "Usuario o password incorrecto")

}
func logoutHandler(w http.ResponseWriter, r *http.Request) {

	usuarioActual = ""

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func registroHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		fmt.Fprintf(w, `
		<html>
		<head>
		<title>Registro</title>
		</head>
		<body>

		<h1>Crear cuenta</h1>

		<form method="POST">

		Usuario:<br>
		<input name="usuario"><br><br>

		Password:<br>
		<input type="password" name="password"><br><br>

		<button type="submit">Registrarse</button>

		</form>

		<br>
		<a href="/login">Ir a login</a>

		</body>
		</html>
		`)

		return
	}

	// POST

	user := r.FormValue("usuario")
	pass := r.FormValue("password")

	if user == "" || pass == "" {
		fmt.Fprintf(w, "❌ Completá todos los campos")
		return
	}

	err := crearUsuarioDB(user, pass)

	if err != nil {

		fmt.Fprintf(w, "❌ El usuario ya existe")
		return
	}

	fmt.Fprintf(w, "✅ Usuario creado! <a href='/login'>Login</a>")
}
