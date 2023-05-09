// Instituto Tecnológico y de Estudios Superiores de Monterrey
//
// Implementación de Métodos Computacionales (Grupo 605)
//
// Resaltador de sintaxis paralelo
//
// Juan Carlos Martínez Zacarías A01612967
//
// Este programa es un resaltador de sintaxis del lenguaje de programación Python.
// Recibe como entrada un directorio con archivos de texto (.txt) con código escrito en dicho lenguaje.
// La salida del programa es un directorio con los archivos HTML correspondientes a los archivos de 
// entrada, los cuales cuentan con el resaltado de sintaxis a partir del uso de CSS en el HTML.
// El programa requiere de dos parámetros: el primero es el directorio fuente y el segundo es el 
// directorio de salida.
// Primeramente se realiza una ejecución secuencial, se obtiene su tiempo de ejecución y este se
// imprime en la terminal. Posteriormente, se realiza la ejecución paralela y, de la misma forma,
// el tiempo de ejecución se calcula e imprime en la terminal de usuario. Finalmente, se imprime
// un mensaje que indica una ejecución exitosa y finalizada.
//
// De esta manera, algunos resultados obtenidos fueron los siguientes (para ambas pruebas se utilizaron los mismos
// archivos):
// - Con 2 archivos:
//   - Tiempo de ejecución secuencial: prueba 1 - 104.96 ms, prueba 2 - 124.58 ms
//   - Tiempo de ejecución paralelo: prueba 1 - 55.96 ms, prueba 2 - 69.86 ms
// - Con 3 archivos:
//   - Tiempo de ejecución secuencial: prueba 1 - 272.36 ms, prueba 2 - 207.44 ms 
//   - Tiempo de ejecución paralelo: prueba 1 - 197.62 ms, prueba 2 - 137.77 ms
// - Con 4 archivos:
//   - Tiempo de ejecución secuencial: prueba 1 - 262.3 ms, prueba 2 - 244.54 ms
//   - Tiempo de ejecución paralelo: prueba 1 - 173.68 ms, prueba 2 - 165.59 ms
// - Con 5 archivos:
//   - Tiempo de ejecución secuencial: prueba 1 - 336.86 ms, prueba 2 - 304.60 ms
//   - Tiempo de ejecución paralelo: prueba 1 - 188.10 ms, prueba 2 - 177.92 ms
//
// Por otro lado, al hablar de la complejidad computacional del algoritmo implementado, se trata de un loop
// principal que itera tantas veces como caractéres haya en el archivo de entrada (hablando de una ejecución individual
// y considerando que un caractér puede contar con más de un byte de almacenamiento). Asimismo, cada iteración
// no va más allá de un orden de complejidad constante, puesto que se trata de la lectura de un caractér que es
// sometido a una cantidad fija de comparaciones (por que la cantidad de comparaciones realizadas para el
// elemento en turno siempre es menor o igual al máximo de comparaciones posibles, el cual no depende de la cantidad
// de caractéres dentro del archivo de entrada) y, tras estas comparaciones, se accede a la matriz de transiciones,
// siendo un accesso inmediato e igualmente constante, de modo que las operaciones que ocurran a continuación ya
// formarán parte de la lectura de un nuevo caractér. De esta forma es como la complejidad computacional del 
// algoritmo es una complejidad lineal O(n), donde n es la cantidad de caractéres contenidos en el archivo de
// entrada, por lo que una ejecución en paralelo del programa, en teoría, sigue siendo O(n) al realizarse el mismo
// procedimiento varias veces, pero al mismo tiempo. No obstante, también existen factores que ponen en duda esta
// afirmación, puesto que hay múltiples procesos que hacen uso de los recursos del equipo y que son indispensables para
// su funcionamiento, por lo que la disponibilidad de memoria y núcleos está variando constantemente. De igual manera,
// el procesamiento de un archivo complejo puede exigir más recursos y, al ser compartidos con el resto de procesos
// dada la ejecución paralela, este proceso puede verse afectado negativamente, o bien puede ser indistinto si es que
// el equipo cuenta con bastante capacidad de procesamiento y memoria.
// 
// Así, es posible ver en el caso de prueba de dos archivos la complejidad computacional establecida, puesto
// que la ejecución secuencial toma prácticamente el doble de tiempo que la paralela. No obstante, para los casos de 
// tres, cuatro y cinco archivos, el tiempo de ejecución paralelo es de un poco más de la mitad del secuencial,
// el cual, en teoría, debería de equivaler a la tercera, cuara, y quinta parte de este último, respectivamente, por lo que 
// podría tratarse de esta compartición de recursos que afecta negativamente al procesamiento, si bien hay una
// mejoría que no es la esperada, pero que sí es notable.
//
// Finalmente, es posible hablar del procesamiento de 5 archivos, siendo uno de ellos ciertamente complejo:
// - Ejecución secuencial: prueba 1 - 10.09 s, prueba 2 - 9.19 s
// - Ejecución paralela: prueba 1 - 9.10 s, prueba 1- 9.08 s
//
// Es notable la complejidad de dicho código considerando que ahora la escala de medición pasó de milisegundos
// a segundos. De la misma forma, es muy probable encontrar en este caso de prueba una compartición de recursos
// que merma el procesamiento, a pesar de que hay una mejora de un segundo en la ejecución en paralelo, por lo que
// se concluye que existen estos dos extremos del uso de la programación paralela: uno en el que la repartición
// de recursos a múltiples procesos es indistinta, y otro en el que afecta negativamente el procesamiento,
// lo cual también se observa a un nivel mayor cuando se evalúa el desempeño de una computadora que ejecuta
// un solo programa o una cantidad considerable al mismo tiempo. 

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"log"
	"os"
	"time"
	"sync"
	"unicode/utf8"
)

func main() {
	var wg sync.WaitGroup // Variable para la espera de la ejecución paralela

	// Lectura de parámetros de línea de comando
    folder := os.Args[1]
	outputFolder := os.Args[2]

	// Lectura de archivos del directorio indicado
	files, err := ioutil.ReadDir(folder)
    if err != nil {
        log.Fatal(err)
    }
	filesNum := len(files) // Conteo de archivos leídos

	transitionMatrix := matrixInit() // Importación de la matriz de transiciones

	// Ejecución secuencial del resaltador de sintaxis
	wg.Add(filesNum)
	start := time.Now()
	for i := 0; i < filesNum; i++ {
		syntaxHighlighter(folder + files[i].Name(), outputFolder + "A01612967_squential_" + files[i].Name() + ".html", &transitionMatrix, &wg)
	}
	wg.Wait()
	timeElapsed := time.Since(start)
	fmt.Println("The sequential excecution took ", timeElapsed)

	// Ejecución paralela del resaltador de sintaxis
	wg.Add(filesNum)
	start = time.Now()
	for i := 0; i < filesNum; i++ {
		go syntaxHighlighter(folder + files[i].Name(), outputFolder + "A01612967_parallel_" + files[i].Name() + ".html", &transitionMatrix, &wg)
	}
	wg.Wait()
	timeElapsed = time.Since(start)
	fmt.Println("The parallel excecution took ", timeElapsed)

	// Mensaje final que indica la finalización del programa
    fmt.Println("Done!")
}

// Función syntaxHighlighter. Procesa el archivo de entrada y crea su respectivo archivo html
// Parámetro filePath: dirección del archivo a procesar
// Parámetro outputPath: dirección del archivo de salida
// Parámetro MT: Matriz de transición del DFA que modela la función. Se pasa por referencia
// Parámetro wg: Permite controlar el número de iteraciones de la función
func syntaxHighlighter(filePath string, outputPath string, MT *[155][55]int, wg *sync.WaitGroup) {
    
	defer wg.Done() // Reducción del número de ejecuciones restantes
	
	var idx int64 // Variable para determinar la posición en la que el nuevo contenido se escribirá en el archivo de salida
	
	// Creación del archivo de salida
	f, err := os.Create(outputPath)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

	// Lectura del archivo de salida para su escritura
	outputfilebuffer, err := ioutil.ReadFile(outputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	// Configuración CSS inicial del archivo de salida
	writingString := "<p><font face=\"Arial\" size=\"4px\">\n"
	stringToBytes := []byte(writingString)
	_, err2 := f.Write(stringToBytes)
    if err2 != nil {
        log.Fatal(err2)
    }

	state := 0 // Número de estado en el DFA
    lexeme := "" // Lexema actual identificado
	fileSize := 0 // Tamaño del archivo en bytes
	counter := 1 // Contador auxiliar
	read := true // Indica si es necesario leer un nuevo caractér del archivo de entrada
	c := "" // Caractér actual

	// Lectura del archivo de entrada
	filebuffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
    inputdata := string(filebuffer)

	fileSize = utf8.RuneCountInString(inputdata) // Conteo de caractéres contenidos en el archivo de entrada

	// Creación del escáner para la lectura de cada caractér del archivo de entrada
	data := bufio.NewScanner(strings.NewReader(inputdata))
	data.Split(bufio.ScanRunes)

	// Loop principal. Se ejecutará tantas veces como caractéres haya
	for counter <= fileSize + 1 {
		for state < 200 { // Dentro del DFA, mientras el estado no sea de Aceptación o Rechazo (no identificado)
			if read {
				if counter == fileSize + 1 {
					c = ""
					counter++
				} else {
					for data.Scan() {
						c = data.Text()
						counter++
						break
					}
				}
			} else {
				read = true
			}
			state = MT[state][filter(c)]
			if state < 200 && state != 0 {
				lexeme += c
			}
		}
		if state == 200 {
			// El siguiente caractér ha sido leído
			read = false 
			// Contenido a escribir en el archivo de salida
			writingString = "<div style=\"color:Indigo;display:inline;\">" + lexeme + "</div>" 
			// Contenido en bytes a escribir en el archivo de salida
			stringToBytes = []byte(writingString) 
			// Lectura actual del archivo de salida
			outputfilebuffer, err = ioutil.ReadFile(outputPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// Cálculo de la posición en la que el nuevo contenido será escrito
			idx = int64(len(outputfilebuffer))
			// Escritura del nuevo contenido en la posición obtenida
			_, err3 := f.WriteAt(stringToBytes, idx)
			if err3 != nil {
				log.Fatal(err3)
			}
		} else if state == 201 {
			read = false
			writingString = "<div style=\"color:DarkViolet;display:inline;\">" + lexeme + "</div>"
			stringToBytes = []byte(writingString)
			outputfilebuffer, err = ioutil.ReadFile(outputPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			idx = int64(len(outputfilebuffer))
			_, err3 := f.WriteAt(stringToBytes, idx)
			if err3 != nil {
				log.Fatal(err3)
			}
		} else if state == 202 {
			read = false
			if lexeme == "true" || lexeme == "false" || lexeme == "none" || lexeme == "Elif" || lexeme == "From" || lexeme == "Else" || lexeme == "Except" || lexeme == "Finally" || lexeme == "For" || lexeme == "Nonlocal" || lexeme == "Not" || lexeme =="Try" {
				writingString = "<div style=\"color:Indigo;display:inline;\">" + lexeme + "</div>"
			    stringToBytes = []byte(writingString)
			    outputfilebuffer, err = ioutil.ReadFile(outputPath)
			    if err != nil {
				    fmt.Println(err)
				    os.Exit(1)
			    }
			    idx = int64(len(outputfilebuffer))
			    _, err3 := f.WriteAt(stringToBytes, idx)
			    if err3 != nil {
				    log.Fatal(err3)
			    }
			} else if lexeme == "def" {
				writingString = "<div style=\"color:Maroon;display:inline;\">" + lexeme + "</div><b>"
			    stringToBytes = []byte(writingString)
			    outputfilebuffer, err = ioutil.ReadFile(outputPath)
			    if err != nil {
				    fmt.Println(err)
				    os.Exit(1)
			    }
			    idx = int64(len(outputfilebuffer))
			    _, err3 := f.WriteAt(stringToBytes, idx)
			    if err3 != nil {
				    log.Fatal(err3)
			    }
			} else {
				writingString = "<div style=\"color:Maroon;display:inline;\">" + lexeme + "</div>"
			    stringToBytes = []byte(writingString)
			    outputfilebuffer, err = ioutil.ReadFile(outputPath)
			    if err != nil {
				    fmt.Println(err)
				    os.Exit(1)
			    }
			    idx = int64(len(outputfilebuffer))
			    _, err3 := f.WriteAt(stringToBytes, idx)
			    if err3 != nil {
				    log.Fatal(err3)
			    }
			}
		} else if state == 203 {
			read = false
			writingString = "<div style=\"color:Gold;display:inline;\">" + lexeme + "</div>"
			stringToBytes = []byte(writingString)
			outputfilebuffer, err = ioutil.ReadFile(outputPath)
			if err != nil {
			    fmt.Println(err)
			    os.Exit(1)
		    }
		    idx = int64(len(outputfilebuffer))
		    _, err3 := f.WriteAt(stringToBytes, idx)
		    if err3 != nil {
			    log.Fatal(err3)
		    }
		} else if state == 204 {
			read = false
			if lexeme == "(" {
				writingString = "<div style=\"color:HotPink;display:inline;\"></b>" + lexeme + "</div>"
			    stringToBytes = []byte(writingString)
			    outputfilebuffer, err = ioutil.ReadFile(outputPath)
			    if err != nil {
			        fmt.Println(err)
			        os.Exit(1)
		        }
		        idx = int64(len(outputfilebuffer))
		        _, err3 := f.WriteAt(stringToBytes, idx)
		        if err3 != nil {
			        log.Fatal(err3)
		        }
			} else {
				writingString = "<div style=\"color:HotPink;display:inline;\">" + lexeme + "</div>"
			    stringToBytes = []byte(writingString)
			    outputfilebuffer, err = ioutil.ReadFile(outputPath)
			    if err != nil {
			        fmt.Println(err)
			        os.Exit(1)
		        }
		        idx = int64(len(outputfilebuffer))
		        _, err3 := f.WriteAt(stringToBytes, idx)
		        if err3 != nil {
			        log.Fatal(err3)
		        }
			}
		} else if state == 205 {
			// El último caractér forma parte del lexema
			lexeme += c
			writingString = "<div style=\"color:Lime;display:inline;\">" + lexeme + "</div>"
			stringToBytes = []byte(writingString)
			outputfilebuffer, err = ioutil.ReadFile(outputPath)
			if err != nil {
			    fmt.Println(err)
			    os.Exit(1)
		    }
		    idx = int64(len(outputfilebuffer))
		    _, err3 := f.WriteAt(stringToBytes, idx)
		    if err3 != nil {
			    log.Fatal(err3)
		    }
		} else if state == 206 {
			read = false
			writingString = "<div style=\"color:Gray;display:inline;\">" + lexeme + "</div>"
			stringToBytes = []byte(writingString)
			outputfilebuffer, err = ioutil.ReadFile(outputPath)
			if err != nil {
			    fmt.Println(err)
			    os.Exit(1)
		    }
		    idx = int64(len(outputfilebuffer))
		    _, err3 := f.WriteAt(stringToBytes, idx)
		    if err3 != nil {
			    log.Fatal(err3)
		    }
		} else if state == 207 {
			read = false
			writingString = "<div style=\"color:Red;display:inline;\">" + lexeme + "</div>"
			stringToBytes = []byte(writingString)
			outputfilebuffer, err = ioutil.ReadFile(outputPath)
			if err != nil {
			    fmt.Println(err)
			    os.Exit(1)
		    }
		    idx = int64(len(outputfilebuffer))
		    _, err3 := f.WriteAt(stringToBytes, idx)
		    if err3 != nil {
			    log.Fatal(err3)
		    }
		} else if state == 209 {
			writingString = "<br>"
			stringToBytes = []byte(writingString)
			outputfilebuffer, err = ioutil.ReadFile(outputPath)
			if err != nil {
			    fmt.Println(err)
			    os.Exit(1)
		    }
		    idx = int64(len(outputfilebuffer))
		    _, err3 := f.WriteAt(stringToBytes, idx)
		    if err3 != nil {
			    log.Fatal(err3)
		    }
		} else if state == 210 {
			writingString = "&nbsp;"
			stringToBytes = []byte(writingString)
			outputfilebuffer, err = ioutil.ReadFile(outputPath)
			if err != nil {
			    fmt.Println(err)
			    os.Exit(1)
		    }
		    idx = int64(len(outputfilebuffer))
		    _, err3 := f.WriteAt(stringToBytes, idx)
		    if err3 != nil {
			    log.Fatal(err3)
		    }
		} else if state == 300 {
			lexeme += c
			read = true
			writingString = "<div style=\"display:inline;\">" + lexeme + "</div>"
			stringToBytes = []byte(writingString)
			outputfilebuffer, err = ioutil.ReadFile(outputPath)
			if err != nil {
			    fmt.Println(err)
			    os.Exit(1)
		    }
		    idx = int64(len(outputfilebuffer))
		    _, err3 := f.WriteAt(stringToBytes, idx)
		    if err3 != nil {
			    log.Fatal(err3)
		    }
		}
		if state == 208 {
			// Fin de la entrada. Se cierra la etiqueta que fue abierta al inicio del archivo de salida y se termina el loop principal
			writingString = "\n</font></p>"
			stringToBytes = []byte(writingString)
			outputfilebuffer, err = ioutil.ReadFile(outputPath)
			if err != nil {
			    fmt.Println(err)
			    os.Exit(1)
		    }
		    idx = int64(len(outputfilebuffer))
		    _, err3 := f.WriteAt(stringToBytes, idx)
		    if err3 != nil {
			    log.Fatal(err3)
		    }
			break
		}
		// Una vez obtenido el lexema, tanto este como el estado actual se reinician
        lexeme = ""
        state = 0
	}
}

// Función filter. Obtiene el valor correspondiente al carácter leído como entrada en la matriz de transición
// Parámetro c: Caractér leído
// Valor de retorno entry: Valor numérico correspondiente al caractér dentro de la matriz de transición
func filter(c string)(entry int) {
	if c == "a" {
		entry = 0
		return
	} else if c == "b" {
		entry = 1
		return
	} else if c == "c" {
		entry = 2
		return
	} else if c == "d" {
		entry = 3
		return
	} else if c == "e" || c == "E" {
		entry = 4
		return
	} else if c == "f" || c == "F" {
		entry = 5
		return
	} else if c == "g" {
		entry = 6
		return
	} else if c == "h" {
		entry = 7
		return
	} else if c == "i" {
		entry = 8
		return	
	} else if c == "j" {
		entry = 9
		return
	} else if c == "k" {
		entry = 10
		return
	} else if c == "l" {
		entry = 11
		return
	} else if c == "m" {
		entry = 12
		return
	} else if c == "n" || c == "N" {
		entry = 13
		return
	} else if c == "o" {
		entry = 14
		return
	} else if c == "p" {
		entry = 15
		return
	} else if c == "q" {
		entry = 16
		return
	} else if c == "r" {
		entry = 17
		return
	} else if c == "s" {
		entry = 18
		return
	} else if c == "t" || c == "T" {
		entry = 19
		return
	} else if c == "u" {
		entry = 20
		return
	} else if c == "v" {
		entry = 21
		return
	} else if c == "w" {
		entry = 22
		return
	} else if c == "x" {
		entry = 23
		return
	} else if c == "y" {
		entry = 24
		return
	} else if c == "z" {
		entry = 25
		return
	} else if (c >= "a" && c <= "z") || (c >= "A" && c <= "Z") || (c >= "À" && c <= "ÿ") {
		entry = 1
		return
	} else if c == "_" {
		entry = 26
		return
	} else if c == "~" {
		entry = 27
		return
	} else if c == "(" || c == ")" || c == "[" || c == "]" || c == "{" || c == "}" || c == ":" || c == "," || c == ";" || c == "`" {
		entry = 28
		return
	} else if c == "'" {
		entry = 29
		return
	} else if c == "#" {
		entry = 30
		return
	} else if c >= "0" && c <= "9" {
		entry = 31
		return
	} else if c == "\n" {
		entry = 32
		return
	} else if c == " " || c == "\t" {
		entry = 33
		return
	} else if c == "*" {
		entry = 34
		return
	} else if c == "/" {
		entry = 35
		return
	} else if c == "<" {
		entry = 36
		return
	} else if c == ">" {
		entry = 37
		return
	} else if c == "=" {
		entry = 38
		return
	} else if c == "." {
		entry = 39
		return
	} else if c == "+" {
		entry = 40
		return
	} else if c == "-" {
		entry = 41
		return
	} else if c == "%" {
		entry = 42
		return
	} else if c == "&" {
		entry = 43
		return
	} else if c == "|" {
		entry = 44
		return
	} else if c == "^" {
		entry = 45
		return
	} else if c == "!" {
		entry = 46
		return
	} else if c == "\"" {
		entry = 47
		return
	} else if c == "" {
		entry = 48
		return
	} else if c == "\\" {
		entry = 50
		return
	} else {
		entry = 49
		return
	}
}

// Función matrixInit. Crea una instancia de la matriz de transiciones que modela el resaltador de sintaxis.
//Valor de retorno matrix: matriz de transiciones del DFA que modela el programa
//
//Tokens
// 200 - Identificador
// 201 - Identificador privado
// 202 - Palabra reservada
// 203 - Operadores
// 204 - Delimitadores
// 205 - Cadenas
// 206 - Comentarios
// 207 - Números
// 209 - Salto de línea
// 210 - Espacios en blanco
// 208 - Fin de la entrada
// 300 - No identificado
//
// Matriz de transición: Codificación de DFA
// [fila, columna] = [estado no final, transición]
// Estados > 199 son finales
// Caso especial: estado 300 - No identificado
//
func matrixInit()(matrix [155][55]int) {
	//      a    b    c    d    e    f    g    h    i    j    k    l    m    n    o    p    q    r    s    t    u    v    w    x    y    z    _    ~  del    '    #  dig   \n   ' '\t    *    /    <    >    =    .    +    -    %    &    |    ^    !    "   ''  otro   \
	MT := [155][55]int {
		{   3,   4,   5,   6,   7,   8,   9,   1,  10,   1,   1,  11,   1,  12,  13,  14,   1,  15,   1,  16,   1,   1,  17,   1,  18,   1,   2,  19,  20,  21,  22,  23, 209,    210,  47,  48,  49,  50,  51,  20,  53,  54,  55,  56,  57,  58,  59,  60, 208,  300,  20},
		{   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
		{   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2,   2, 201, 201, 201, 201,   2, 201,    201, 201, 201, 201, 201, 201, 201, 201, 201, 201, 201, 201, 201, 201, 201, 201,  201, 201},
		{   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  24,   1,   1,   1,   1,  25,   1,   1,   1, 135,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  26,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  27,   1,   1,  28,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,  29,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  30,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  31,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 143,   1,   1,   1,   1,   1,   1,   1,  32,   1,   1,   1,   1,   1,  33,   1,   1,  34,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  35,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,  36, 128,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {  37,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  38,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {  39,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {  41,   1,   1,   1,  42,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  43,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,  44,  45,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,  46,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        {  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21, 205,  21,  21,  21,     21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,  21,   21, 147},
        {  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22, 206,     22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22,  22, 206,   22,  22},
        { 207, 207, 207, 207, 129, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207,  23, 207,    207, 207, 207, 207, 207, 207, 132, 207, 207, 207, 207, 207, 207, 207, 207, 207,  207, 207},
        {   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  61,   1,   1,   1,   1,   1, 133,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,  62,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {  63,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  64,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,  65,   1,   1,   1,   1,   1,   1,   1,   1,   1,  66,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,  67,   1,  68,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  69,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  70,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  71,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  72,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  73,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 138,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  74,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,  75,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,  76,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  77,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 146,   1,   1,   1, 128,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,  78,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,  79,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,  80,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203,  81, 203, 203, 203,  82, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203,  83, 203, 203,  84, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203,  85, 203,  86, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203,  87,  88, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204,  89, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        {}, 
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203,  91, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203,  92, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203,  94, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203,  95, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203,  96, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203,  97, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203,  98, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        {  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,     60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60,  60, 205,  60,   60, 148},
        {   1,   1,   1,   1,  99,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 100,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 101,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 102,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1, 103,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 104,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1, 105,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 106,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1, 107,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 108,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 109,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 110,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 111,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 112,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203, 113, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203, 114, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203, 115, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203, 116, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        {},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        {},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        { 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,    203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203, 203,  203, 203},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 117,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1, 118,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 119,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 120,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 121,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 122,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1, 123,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 124,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        { 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,    204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204, 204,  204, 204},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 125,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 126,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 127,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 202, 202, 202, 202,   1, 202,    202, 202, 202, 202, 202, 202, 202, 202, 202, 202, 202, 202, 202, 202, 202, 202,  202, 202},
        { 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 130, 300,    300, 300, 300, 300, 300, 300, 300, 131, 131, 300, 300, 300, 300, 300, 300, 300,  300, 300},
        { 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 130, 207,    207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207,  207, 207},
        { 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 130, 300,    300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,  300, 300},
        { 207, 207, 207, 207, 129, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 132, 207,    207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207, 207,  207, 207},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 134,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 136,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1, 137,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1, 139,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 140,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1, 141,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        { 142,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},    
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 144,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 145,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},
        {   1,   1,   1,   1, 128,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1,   1, 200, 200, 200, 200,   1, 200,    200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200,  200, 200},      
        { 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 149, 150, 150, 150,    150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150,  150, 150},
        { 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152,    152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 151, 152,  152, 152},
        { 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 205, 149, 149, 149,    149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149,  149, 149},
        { 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 205, 150, 150, 150,    150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150, 150,  150, 150},
        { 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151,    151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 151, 205, 151,  151, 151},
        { 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152,    152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 152, 205, 152,  152, 152},
	}
	matrix = MT
	return
}