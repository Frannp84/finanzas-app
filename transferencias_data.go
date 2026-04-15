package main

import (
    "encoding/json"
    "os"
)

func guardarTransferencias() {

    file, _ := os.Create("transferencias.json")
    defer file.Close()

    json.NewEncoder(file).Encode(transferencias)

}

func cargarTransferencias() {

    file, err := os.Open("transferencias.json")

    if err != nil {
        return
    }

    defer file.Close()

    json.NewDecoder(file).Decode(&transferencias)

}