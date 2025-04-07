/* Este programa genera un Dockerfile y un archivo JSON de configuración
   a partir de la entrada del usuario. Luego, construye una imagen Docker
   y ejecuta un contenedor basado en esa imagen.
   Autor: Walter Avendaño
   Fecha: 04/04/2025
(*/

package main

//Cargando librerias necesarias
import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Definición de la estructura para el JSON
type DockerConfig struct {
	ID        string `json:"id"`
	BaseImage string `json:"baseImage"`
	//	BaseLabel  string `json:"baseLabel"`
	WorkDir    string `json:"workDir"`
	CopyFiles  string `json:"copyFiles"`
	RunCmds    string `json:"runCmds"`
	ExposePort int    `json:"exposePort"`
	DefaultCmd string `json:"defaultCmd"`
}

// Función principal
// Esta función se encarga de crear un Dockerfile y un archivo JSON de configuración
func main() {
	// Crea un lector de entrada
	reader := bufio.NewReader(os.Stdin)

	// Pide los parámetros al usuario
	fmt.Print("Ingrese el ID del Dockerfile: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print("Ingrese la imagen base (por ejemplo, node:8): ")
	baseImage, _ := reader.ReadString('\n')
	baseImage = strings.TrimSpace(baseImage)

	/*	color.New(color.BgBlue).Println("Ingrese la etiqueta a usar (por ejemplo, Node-app): ")
		baseLabel, _ := reader.ReadString('\n')
		baseLabel = strings.TrimSpace(baseLabel)
	*/
	fmt.Print("Ingrese el directorio de trabajo (por ejemplo, /usr/src/app): ")
	workDir, _ := reader.ReadString('\n')
	workDir = strings.TrimSpace(workDir)

	fmt.Print("Ingrese los archivos a copiar (por ejemplo, .): ")
	copyFiles, _ := reader.ReadString('\n')
	copyFiles = strings.TrimSpace(copyFiles)

	fmt.Print("Ingrese los comandos a ejecutar (por ejemplo, npm install): ")
	runCmds, _ := reader.ReadString('\n')
	runCmds = strings.TrimSpace(runCmds)

	fmt.Print("Ingrese el puerto a exponer (por ejemplo, 3000): ")
	exposePortStr, _ := reader.ReadString('\n')
	exposePortStr = strings.TrimSpace(exposePortStr)
	exposePort, err := strconv.Atoi(exposePortStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Ingrese el comando por defecto (por ejemplo, node index.js): ")
	defaultCmd, _ := reader.ReadString('\n')
	defaultCmd = strings.TrimSpace(defaultCmd)

	// Crea la estructura de datos para el JSON
	config := DockerConfig{
		ID:        id,
		BaseImage: baseImage,
		//	BaseLabel:  baseLabel,
		WorkDir:    workDir,
		CopyFiles:  copyFiles,
		RunCmds:    runCmds,
		ExposePort: exposePort,
		DefaultCmd: defaultCmd,
	}

	// Genera el contenido del Dockerfile
	dockerfileContent := fmt.Sprintf(`FROM %s
WORKDIR %s
COPY %s .
RUN %s
EXPOSE %d
CMD ["%s"]
`, config.BaseImage, config.WorkDir, config.CopyFiles, config.RunCmds, config.ExposePort, config.DefaultCmd)

	// Escribe el contenido en un archivo Dockerfile
	err = os.WriteFile("Dockerfile", []byte(dockerfileContent), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Dockerfile generado exitosamente")

	// Serializa la estructura a JSON
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Escribe el JSON a un archivo
	err = os.WriteFile("docker_config.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Archivo JSON generado exitosamente")

	// Construye la imagen Docker
	buildCmd := exec.Command("docker", "build", "-t", fmt.Sprintf("%s:%s", config.ID, "latest"), ".")
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
	fmt.Println("Imagen Docker construida exitosamente")

	// Crea y ejecuta el contenedor
	runCmd := exec.Command("docker", "run", "-d", "-p", fmt.Sprintf("%d:%d", config.ExposePort, config.ExposePort), fmt.Sprintf("%s:%s", config.ID, "latest"))
	output, err = runCmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
	fmt.Println("Contenedor creado y ejecutado exitosamente")
}
