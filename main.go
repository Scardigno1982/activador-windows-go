package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
)

func runCommand(command string) (string, error) {
	cmd := exec.Command("powershell", "-Command", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func findEdition(osEdition string, licenses map[string]string) (string, bool) {
	for key, value := range licenses {
		if strings.Contains(strings.ToLower(osEdition), strings.ToLower(key)) {
			return value, true
		}
	}
	return "", false
}

func main() {
	// Verificar si se está ejecutando en Windows
	if runtime.GOOS != "windows" {
		fmt.Println("Este programa está diseñado para ejecutarse en Windows.")
		return
	}

	// Verificar si se está ejecutando como administrador
	cmd := "([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)"
	output, err := runCommand(cmd)
	if err != nil || strings.TrimSpace(output) != "True" {
		fmt.Println("Debes ejecutar este programa como administrador.")
		return
	}

	// Obtener la edición del sistema operativo
	cmd = "(Get-CimInstance Win32_OperatingSystem).Caption"
	output, err = runCommand(cmd)
	if err != nil {
		fmt.Println("Error al obtener la edición del sistema operativo.")
		return
	}
	osEdition := strings.TrimSpace(output)

	// Ruta al archivo JSON de licencias
	jsonFile := "licencias.json"

	// Leer el contenido del archivo JSON
	fileContent, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		fmt.Println("Error al leer el archivo JSON de licencias.")
		return
	}

	// Crear un mapa para almacenar las licencias
	var licenses map[string]string
	err = json.Unmarshal(fileContent, &licenses)
	if err != nil {
		fmt.Println("Error al analizar el contenido del archivo JSON.")
		return
	}

	// Buscar la licencia correspondiente
	licenseKey, found := findEdition(osEdition, licenses)

	// Mostrar la licencia
	if found {
		fmt.Printf("La licencia para la edición %s es: %s\n", osEdition, licenseKey)
	} else {
		fmt.Printf("No se encontró una licencia para la edición %s en el archivo JSON.\n", osEdition)
	}
}
