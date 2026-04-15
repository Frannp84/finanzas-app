package main

import (
	"encoding/json"
	"os"
)

type Usuario struct {
	Usuario  string `json:"usuario"`
	Password string `json:"password"`
}

var usuarios []Usuario

func cargarUsuarios() {

	file, err := os.Open("usuarios.json")

	if err != nil {

		return

	}

	defer file.Close()

	json.NewDecoder(file).Decode(&usuarios)

}

func validarUsuario(user string, pass string) bool {

	for _, u := range usuarios {

		if u.Usuario == user && u.Password == pass {

			return true

		}

	}

	return false

}
