package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("ERROR leyendo usuarios.json:", err)
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
