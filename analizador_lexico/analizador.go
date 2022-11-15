/*
	# def: Analizador léxico para el lenguaje de consultas estructurado (SQL)
	# Tabla de símbolos: Inicialmente contendrá cargadas las palabras resevardas y adicionalmente la utilizaremos para almacenar las variables definidas por el usuario
		- Estructura para palabras reservadas: Nombre, Tipo, Bytes, NE, NE, SQL
			* NE: NO ESPECIFICADO
		- Estructura para variables de usuario: Nombre, Tipo, Bytes, Linea(declarada), Lineas(Uso), Definida por
*/

package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Definición de HashTable

/*
#tamHashTable: Representa el tamaño del arreglo utilizado como hashtable, el tamaño es representado por la cantidad de palabras reservadas en SQL * 2 (únicamente las que fueron abarcadas)
*/
const tamHashTable = 150

type hashtable struct {
	array [tamHashTable]*bucket
}

type bucket struct {
	head *bucketNode
}

type bucketNode struct {
	key  [1][5]string
	next *bucketNode
}

//Funciones principales para manipulación del HashTable

func hash(key string) int {
	indice := 0
	for _, letras := range key {
		indice += int(letras)
	}
	return indice % tamHashTable
}

func (h *hashtable) insertInHashTable(key [1][5]string) {
	index := hash(key[0][0])
	h.array[index].insertBucketNode(key)
}

func (b *bucket) insertBucketNode(k [1][5]string) {
	keyNode := &bucketNode{key: k}
	keyNode.next = b.head
	b.head = keyNode
}

func (h *hashtable) deleteInHashTable(key string) {
	index := hash(key)
	h.array[index].deleteBucketNode(key)
}

func (b *bucket) deleteBucketNode(k string) {
	var nodoAnterior *bucketNode
	nodoSearch := b.head
	if nodoSearch != nil {
		for nodoSearch != nil {
			if nodoSearch.key[0][0] != k {
				nodoAnterior = nodoSearch
				nodoSearch = nodoSearch.next
			} else {
				if nodoAnterior != nil {
					nodoAnterior.next = nodoSearch.next
				} else {
					nodoAnterior = nodoSearch.next
				}
				b.head = nodoAnterior
				break
			}
		}
	} else {
		llenarLog("Key no encontrada")
	}

}

func (h *hashtable) searchInHashTable(key string) bool {
	index := hash(key)
	return h.array[index].searchBucketNode(key)
}

func (b *bucket) searchBucketNode(k string) bool {
	nodoActual := b.head
	for nodoActual != nil {
		if nodoActual.key[0][0] == k {
			return true
		}
		nodoActual = nodoActual.next
	}
	return false
}

func (h *hashtable) searchTipoInHashTable(key string) string {
	index := hash(key)
	return h.array[index].searchTipoBucketNode(key)
}

func (b *bucket) searchTipoBucketNode(k string) string {
	var retorno string
	nodoActual := b.head
	for nodoActual != nil {
		if nodoActual.key[0][0] == k {
			retorno = nodoActual.key[0][1]
		}
		nodoActual = nodoActual.next
	}
	return retorno
}

func (h *hashtable) modificarInHashTable(key string, linea string) {
	index := hash(key)
	h.array[index].modificarBucketNode(key, linea)
}

func (b *bucket) modificarBucketNode(key string, linea string) {
	var nodoAnterior *bucketNode = nil
	nodoModificar := b.head
	if nodoModificar != nil {
		for nodoModificar != nil {
			if nodoModificar.key[0][0] != key {
				nodoAnterior = nodoModificar
				nodoModificar = nodoModificar.next
			} else {
				if nodoModificar.key[0][3] != "NE" {
					nodoModificar.key[0][2] = nodoModificar.key[0][3] + "," + linea
				} else {
					nodoModificar.key[0][3] = linea
				}
				if nodoAnterior != nil {
					nodoAnterior.next = nodoModificar
				} else {
					nodoAnterior = nodoModificar
				}
				b.head = nodoAnterior
				break
			}
		}
	} else {
		llenarLog("Key no encontrada")
	}
}

func (h *hashtable) toString() {
	resultado := &hashtable{}
	var contenido string
	for i := range resultado.array {
		if h.array[i].mostrarContenido() != "" {
			if contenido != "" {
				contenido += "\n" + h.array[i].mostrarContenido()
			} else {
				contenido = h.array[i].mostrarContenido()
			}
		}
	}
	generarArchivos([]byte(contenido), "TS_actualizada.csv")
}

func (b *bucket) mostrarContenido() string {
	var data string
	current := b.head
	if current != nil {
		for current != nil {
			if data != "" {
				data += "\n\t" + current.key[0][0] + "," + current.key[0][1] + "," + current.key[0][2] + "," + current.key[0][3] + "," + current.key[0][4]
			} else {
				data = "\t" + current.key[0][0] + "," + current.key[0][1] + "," + current.key[0][2] + "," + current.key[0][3] + "," + current.key[0][4]
			}
			current = current.next
		}
	}
	return data
}

func initializeHashTable() *hashtable {
	resultado := &hashtable{}
	for i := range resultado.array {
		resultado.array[i] = &bucket{}
	}
	return resultado
}

var tablaSimbolos *hashtable = initializeHashTable()
var bandera bool = false
var posibleComentario string
var numLinea int = 0
var tokensLexemas string
var contTable int
var nuevoElemento [1][5]string
var log string
var guardaTable string
var banderaTable bool

func main() {
	var nombreArchivo string
	llenarTS()
	llenarLog("LOG:\tAnalisis iniciado")
	fmt.Println("Ingrese el nombre del archivo con su extensión")
	fmt.Scanln(&nombreArchivo)
	dataLimpia := AnalizarSQL(nombreArchivo, 1)
	if dataLimpia != "" {
		nombreArchivo = "dataLimpia_" + nombreArchivo
		generarArchivos([]byte(dataLimpia), nombreArchivo)
		llenarLog("LOG:\tObtención de Token, iniciada")
		dataTokensLexemas := AnalizarSQL(nombreArchivo, 2)
		nombreArchivo = strings.Replace(nombreArchivo, "dataLimpia_", "tokensLexemas_", 1)
		nombreArchivo = strings.Replace(nombreArchivo, "sql", "txt", 1)
		dataTokensLexemas = tokenFormat(dataTokensLexemas)
		generarArchivos([]byte(dataTokensLexemas), nombreArchivo)
		tablaSimbolos.toString()
		nombreArchivo = strings.Replace(nombreArchivo, "tokensLexemas_", "Logs_", 1)
		llenarLog("LOG:\tAnalisis finalizado")
		generarArchivos([]byte(log), nombreArchivo)
	} else {
		llenarLog("ERROR:\tRectifica el nombre de tu archivo y vuelve a intentar")
		nombreArchivo = strings.Replace(nombreArchivo, nombreArchivo, "Logs_.txt", 1)
		llenarLog("LOG:\tAnalisis finalizado")
		generarArchivos([]byte(log), nombreArchivo)
	}
}

//Funciones para analisis

/*
# def: Función que nos actualiza la información de la variable que va almacenando los logs
# data (in): Parametro que contiene la información que se irá guardando en la variable de los Logs
*/
func llenarLog(data string) {
	if log != "" {
		log += "\n" + data
	} else {
		log = data
	}
}

/*
# def: Función que sirve únicamente para darle formato para que se vean más ordenados los tokens almacenados
# return: Devuelve un string con la data formateada
# data (in): Parametro que contiene los tokens a los que se les dará formato
*/
func tokenFormat(data string) string {
	var retorno string = fmt.Sprintf("%-25s  %-25s %32s", "Token", "Valor", "Linea")
	if data != "" {
		lineaToken := Tokenizador([]byte(data), "\n")
		for _, linea := range lineaToken {
			palabras := Tokenizador([]byte(linea), ", ")
			if retorno != "" {
				retorno += fmt.Sprintf("\n%-25s  %-25s  %32s", palabras[0], palabras[1], palabras[2])
			} else {
				retorno = fmt.Sprintf("\n%-25s  %-25s  %32s", palabras[0], palabras[1], palabras[2])

			}
		}
	} else {
		llenarLog("ERROR:\tTokens vacíos, verifique sus archivos y vuelva a intentar")
	}
	return retorno
}

/*
# def: Función que nos ayudará a cargar la tabla de símbolos en la primera carga del analizador
# dataTS: array que nos sirve para acomodar la información de palabras propias de SQL que serán guardadas en la tabla de simbolos
# nombreFile: Tiene seteado el nombre del archivo que contiene las palabras propias de SQL
# path: Tiene seteada la ruta del archivo ue contiene las palabras propias de SQL
# ix: Variable de conteo para asignar el orden correctamente de los valores la Tabla de Símbolos
*/
func llenarTS() {
	var ix int = 0
	var dataTS [1][5]string
	var nombreFile string = "file_TS.csv"
	var path string = "../data_sources/" + nombreFile
	file, err := ioutil.ReadFile(path)
	if err == nil {
		renglonData := Tokenizador([]byte(file), "\n")
		for _, renglon := range renglonData {
			ix = 0
			palabrasRenglon := Tokenizador([]byte(renglon), ";")
			for _, palabra := range palabrasRenglon {
				dataTS[0][ix] = string(palabra)
				ix++
			}
			tablaSimbolos.insertInHashTable(dataTS)
		}
	} else {
		llenarLog("ERROR:\tEl archivo: " + nombreFile + ", no se encontró o no se pudo abrir")
	}
}

/*
# def: Función que sirve para hacer los analisis según se requiera a los archivos SQL
return: Retornará información que se genera al final del analisis de las 2 opciones
# nombreSQL (in): Parametro que contendrá el nombre del archivo que se desea analizador
# Op (in): Parametro que contendrá la opción de analisis que deseamos
  - 1: Opción que sirve para analizar el archivo, eliminar comentarios y espacios (hace un analisis rápido de comentarios para identificar errores)
  - 2: Opción que servirá para hacer la extracción de token y añadir, modificar elementos a la tabla de símbolos

# numLinea: La reiniciamos a cero con cada analisis que se realiza, para contar las líneas
# identificador: Variable utilizada en la opción 2, para almacenar palabra por palabra los identificadores (palabra reservada, variable o identificador)
# msj: Variable para almacenar las cadenas de literales
# delim: Variable que almacenará los delimitadores en el código
# dataRetorno: Almacenará la información que se va generando ya sea de la opción 1 o 2
# runaLetra: Variable para convertir las letras en caracteres
*/
func AnalizarSQL(nombreSQL string, Op int) string {
	numLinea = 0
	var identificador string
	var delim string
	var msj string
	var runaLetra rune
	var dataRetorno string
	SQL, err := ioutil.ReadFile("../data_sources/" + nombreSQL)
	if err != nil {
		llenarLog("ERROR:\tEl archivo: " + nombreSQL + ", no se encontró o no se pudo abrir")
	} else {
		lineasSQL := Tokenizador(SQL, "\n")
		for _, renglonSQL := range lineasSQL {
			switch Op {
			case 1:
				if strings.TrimSpace(renglonSQL) != "" {
					analizarSimbolos(renglonSQL)
					switch ExisteComentario(renglonSQL) {
					case true:
						if !bandera {
							AnalizarComentario(renglonSQL, numLinea+1)
							dataPrecomentario := extraerCodigo(renglonSQL)
							if dataPrecomentario != "" {
								dataRetorno = string(fmt.Appendln([]byte(dataRetorno), dataPrecomentario))
							}
						} else {
							AnalizarComentario(renglonSQL, numLinea+1)
						}
						break
					case false:
						if !bandera {
							if posibleComentario != "" {
								dataRetorno = string(fmt.Appendln([]byte(dataRetorno), posibleComentario))
								posibleComentario = ""
							}
							dataRetorno = string(fmt.Appendln([]byte(dataRetorno), renglonSQL))
						} else {
							llenarPosibleComentario(renglonSQL)
						}
						break
					}
				}
				break
			case 2:
				palabrasLinea := Tokenizador([]byte(renglonSQL), " ")
				for _, palabra := range palabrasLinea {
					if strings.TrimSpace(palabra) != "" {
						letras := Tokenizador([]byte(palabra), "")
						for _, letra := range letras {
							runaLetra = rune(letra[0])
							if unicode.IsLetter(runaLetra) || (unicode.IsNumber(runaLetra) && identificador != "") || runaLetra == '_' || runaLetra == '@' {
								if !bandera {
									if runaLetra == '@' && identificador != "" {
										agregarIdentificadores(identificador)
										identificador = ""
									}

									if unicode.IsNumber(runaLetra) {
										if !verificarAgregar(identificador, 1) {
											if identificador != "" {
												identificador += string(runaLetra)
											} else {
												identificador = string(runaLetra)
											}
										} else {
											identificador = ""
										}
									} else {
										if identificador != "" {
											identificador += string(runaLetra)
										} else {
											identificador = string(runaLetra)
										}
									}
								} else {
									if msj != "" {
										msj += string(runaLetra)
									} else {
										msj = string(runaLetra)
									}
								}
							}

							if (unicode.IsGraphic(runaLetra) && !unicode.IsLetter(runaLetra) && !unicode.IsNumber(runaLetra)) && runaLetra != '_' && runaLetra != '@' && runaLetra != '.' {
								if runaLetra != '\'' {
									if identificador != "" {
										agregarIdentificadores(identificador)
										identificador = ""
									}

									if (delim == "") && (runaLetra == '<' || runaLetra == '!' || runaLetra == '>') {
										delim = string(runaLetra)
									} else {
										if (delim == "<" || delim == "!" || delim == ">") && (runaLetra == '=' || runaLetra == '>' || runaLetra == '<') {
											delim += string(runaLetra)
											guardarDelimitador(delim)
											delim = ""
										} else {
											guardarDelimitador(string(runaLetra))
											delim = ""
										}
									}

								} else {
									if !bandera {
										if msj != "" {
											msj += string(runaLetra)
										} else {
											msj = string(runaLetra)
											bandera = true
										}
									} else {
										msj += string(runaLetra)
										insertarTokenLexema(msj, "Cadena de Literales", numLinea+1)
										msj = ""
										bandera = false
									}
								}

							}

						}

						if identificador != "" {
							agregarIdentificadores(identificador)
							identificador = ""
						}

						if delim != "" {
							guardarDelimitador(delim)
							delim = ""
						}
					}

					if bandera {
						msj += " "
					}
				}
				if bandera {
					insertarTokenLexema(strings.TrimSpace(msj), "Cadena de Literales", numLinea+1)
					msj = ""
					bandera = false
				}
				dataRetorno = tokensLexemas
				break
			}
			numLinea++
		}
	}
	llenarLog("LOG:\tSe leyeron un total de: " + strconv.Itoa(numLinea+1) + " líneas de código")
	return strings.TrimSpace(dataRetorno)
}

/*
# def: Función que sirve para confirmar si un delimitador, identificador existe en la tabla de símbolos
# return: Devuelve true si la data existe ya en la tabla de símbolos y fue actualizado el espacio de donde fue usada
# data (in): Parametro que almacena la data que vamos a verificar
*/
func verificarAgregar(data string, op int) bool {
	if op == 1 {
		data = strings.ToUpper(data)
	}

	if tablaSimbolos.searchInHashTable(data) {
		tablaSimbolos.modificarInHashTable(data, strconv.Itoa(numLinea+1))
		tipo := tablaSimbolos.searchTipoInHashTable(data)
		insertarTokenLexema(data, tipo, numLinea+1)
		return true
	}
	return false
}

/*
# def: Función que nos ayudará para identificar errores en parentesis y en '
# data (in): Recibe toda la línea del archivo para hacer las respectivas busquedas de errores
*/
func analizarSimbolos(data string) {

	if strings.Contains(data, "CREATE TABLE") {
		guardaTable = data
		if strings.Contains(data, "(") {
			contTable++
			banderaTable = true
		} else {
			llenarLog("ERROR:\tError en la apertura de parentesis, falta por aperturar parentesis, en la línea: " + strconv.Itoa(numLinea))
			contTable = 0
			banderaTable = false
		}
	} else {
		if strings.Contains(data, "CREATE TABLE") || strings.Contains(data, "DROP") || strings.Contains(data, "ALTER") {
			if banderaTable {
				if strings.Contains(data, ")") {
					contTable = 0
					banderaTable = false
				} else {
					llenarLog("ERROR:\tError en el cierre de parentesis, falta por cerrar parentesis, en la línea: " + strconv.Itoa(numLinea))
					banderaTable = false
					contTable = 0
				}
			}
		} else {
			if strings.Contains(data, ")") {
				contTable = 0
				banderaTable = false
			}
		}
	}

	if guardaTable == "" {
		parentesis := strings.ContainsAny(data, "()")
		if parentesis {
			cAbiertos := strings.Count(data, "(")
			cCerrados := strings.Count(data, ")")
			if cCerrados != cAbiertos {
				if cCerrados < cAbiertos {
					if contTable == 0 {
						llenarLog("ERROR:\tError en la apertura de parentesis, falta por cerrar: " + strconv.Itoa(cAbiertos-cCerrados) + " parentesis, en la línea: " + strconv.Itoa(numLinea+1))
					}
				} else {
					if contTable > 0 {
						llenarLog("ERROR:\tError en el cierre de parentesis, falta por aperturar: " + strconv.Itoa(cCerrados-cAbiertos) + " parentesis, en la línea: " + strconv.Itoa(numLinea+1))
					}
				}
			}
		}
	}

	comillas := strings.Contains(data, "'")
	if comillas {
		cComillas := strings.Count(data, "'")
		if cComillas%2 != 0 {
			llenarLog("ERROR:\tError en la cadena de texto, faltan comillas ('), en la línea: " + strconv.Itoa(numLinea+1))
		}
	}
	guardaTable = ""
}

/*
# def: Función para agregar a la tabla de símbolos los identificadores que se encuentren en el código
# data (in): Parametro que contendrá el identificador
# nuevaDataTS: array que sirve para ordenar la data que será guardada en la tabla de símbolos (cuando se requiera)
# reIdentificador: constante que contiene la expresión regular para validar una variable
# reReservada: constante que contiene la expresión regular para validar una palabra reservada
# reIdentificadorND: constante que contiene la expresión regular para validar un identificador
# coincide: Variable que contiene la validación de si el identificador
*/
func agregarIdentificadores(data string) {
	var nuevaDataTS [1][5]string
	reIdentificador := "^@.*$"
	reReservada := "^[^@]\\D*$"
	reIdentificadorND := "^[^@].*$"
	if data != "" {
		coincide, _ := regexp.Match(reReservada, []byte(data))
		if coincide {
			if !verificarAgregar(data, 1) {
				insertarTokenLexema(data, "Identificador", numLinea+1)
				nuevaDataTS[0][0] = data
				nuevaDataTS[0][1] = "Identificador"
				nuevaDataTS[0][2] = strconv.Itoa(numLinea + 1)
				nuevaDataTS[0][3] = "NE"
				nuevaDataTS[0][4] = "Usuario"
				tablaSimbolos.insertInHashTable(nuevaDataTS)
			}
		} else {
			coincide, _ := regexp.Match(reIdentificadorND, []byte(data))
			if coincide {
				//fmt.Println("data: " + data)
				if !verificarAgregar(data, 0) {
					insertarTokenLexema(data, "Identificador", numLinea+1)
					nuevaDataTS[0][0] = data
					nuevaDataTS[0][1] = "Identificador"
					nuevaDataTS[0][2] = strconv.Itoa(numLinea + 1)
					nuevaDataTS[0][3] = "NE"
					nuevaDataTS[0][4] = "Usuario"
					tablaSimbolos.insertInHashTable(nuevaDataTS)
				}
			} else {
				coincide, _ := regexp.Match(reIdentificador, []byte(data))
				if coincide {
					//fmt.Println("Declaración: " + data)
					if !verificarAgregar(data, 0) {
						insertarTokenLexema(data, "Declaración de variable", numLinea+1)
						nuevaDataTS[0][0] = data
						nuevaDataTS[0][1] = "Declaración de variable"
						nuevaDataTS[0][2] = strconv.Itoa(numLinea + 1)
						nuevaDataTS[0][3] = "NE"
						nuevaDataTS[0][4] = "Usuario"
						tablaSimbolos.insertInHashTable(nuevaDataTS)
					}
				}
			}
		}

	}
}

/*
# def: Función que nos ayuda a guardar un delimitador en la tabla de símbolos
# reSimbolo: constante que contiene la expresión regular para validar un delimitador
# coincide: Variable que contiene la validación de si el delimitador es enrealidad un símbolo
*/
func guardarDelimitador(delim string) {
	reSimbolo := "\\W"
	coincide, _ := regexp.Match(reSimbolo, []byte(delim))
	if coincide {
		//fmt.Println("Declaración: " + identificador)
		if !verificarAgregar(delim, 0) {
			llenarLog("WARNING:\tSimbolo: " + delim + ", en la línea: " + strconv.Itoa(numLinea+1) + ", no identificado dentro de la sintaxís definida para SQL")
			insertarTokenLexema("[W]"+delim, "Delimitador/Separador", numLinea+1)
		}
	}
}

/*
# def: Función para llenar la variable que contendrá los tokens
# return: devuelve los tokens que se han agregado
*/
func insertarTokenLexema(data string, tipo string, num int) string {
	if tokensLexemas != "" {
		tokensLexemas += "\n" + data + ", " + tipo + ", " + strconv.Itoa(num)
	} else {
		tokensLexemas = data + ", " + tipo + ", " + strconv.Itoa(num)
	}
	return tokensLexemas
}

/*
# def: Función que nos ayudará a identificar si la línea tiene comentarios o potenciales comentarios
# return: Devuelve true si la linea contiene potenciales comentarios
# data (in): Parametro para almacenar la línea que deseamos analizar
*/
func ExisteComentario(data string) bool {
	return strings.ContainsAny(data, "-/*")
}

/*
# def: Función para analizar una línea que se ha identificado como potencial comentario
# data (in): Parametro que almacena la línea que parece potencial comentario
# num (in): Parametro que almacema el número de línea donde se identificó el posible comentario
# dataSQL: Quitá los espacios al inicio y fin de la línea
# comentarioSimple: Variable booleana que contendrá true si la linea contiene los simbolos de un comentario simple
# aperturaMulti: Variable booleana que contendrá true si la linea contiene los simbolos de una apertura de comentario multiple
# cierreMulti: Variable booleana que contendrá true si la linea contiene los simbolos de una cierre de comentario multiple
# potencialErr1: Variable booleana que contendrá true si la linea contiene los simbolos de un potencial error
# potencialErr2: Variable booleana que contendrá true si la linea contiene los simbolos de un potencial error
# potencialErr3: Variable booleana que contendrá true si la linea contiene los simbolos de un potencial error
*/
func AnalizarComentario(data string, num int) {
	dataSQL := []byte(strings.TrimSpace(data))
	comentarioSimple := strings.Contains(data, "--")
	aperturaMulti := strings.Contains(data, "/*")
	cierreMulti := strings.Contains(data, "*/")
	potencialErr1 := strings.Contains(data, "-")
	potencialErr2 := strings.Contains(data, "*")
	potencialErr3 := strings.Contains(data, "/")

	if !comentarioSimple {
		if potencialErr1 {
			llenarLog("ERROR:\tError en la sintaxis de comentario, en la línea: " + strconv.Itoa(num))
			bandera = false
		} else {
			if !aperturaMulti {
				if !cierreMulti {
					if potencialErr2 || potencialErr3 {
						if VerificarPotencialError(dataSQL, num) {
							if bandera == true {
								posibleComentario = ""
								bandera = false
							} else {
								bandera = true
							}
						} else {
							llenarPosibleComentario(string(dataSQL))
						}
					}
				} else {
					posibleComentario = ""
					bandera = false
				}
			} else {
				if cierreMulti {
					llenarPosibleComentario(string(dataSQL))
					bandera = false
				} else {
					bandera = true
				}
			}
		}
	} else {
		llenarPosibleComentario(string(dataSQL))
		bandera = false
	}
}

/*
# def: Función que sirve para llenar la variable que contendrá todo el código que se consideraría parte del comentario
# data (in): Contiene la información que se considera parte del comentario
*/
func llenarPosibleComentario(data string) {
	if data != "" {
		if posibleComentario != "" {
			posibleComentario += "\n" + data
		} else {
			posibleComentario = data
		}
	}
}

/*
# def: Función que servirá para verificar los potenciales errores 2 y 3, para evitar que un SELECT * se considere error de comentario
# data (in): Contiene los bytes del texto que se envia para analizar
# num (in): Contiene el número de línea al que pertenece la línea
# palabrasRenglon: Contiene las palabras de la línea tokenizada
*/
func VerificarPotencialError(data []byte, num int) bool {
	palabrasRenglon := Tokenizador(data, " ")
	var resp bool
	if len(palabrasRenglon) > 1 {
		palabrasRenglon[0] = strings.ToUpper(palabrasRenglon[0])
		palabrasRenglon[1] = strings.ToUpper(palabrasRenglon[1])
		palabrasRenglon[2] = strings.ToUpper(palabrasRenglon[2])

		if palabrasRenglon[0] == "SELECT" && (palabrasRenglon[1] == "*" || strings.Contains(palabrasRenglon[1], "COUNT") || strings.Contains(palabrasRenglon[1], "SUM") || strings.Contains(palabrasRenglon[1], "AVG")) && palabrasRenglon[2] == "FROM" {
			if strings.Count(string(data), "*") <= 1 {
				resp = false
			}
		}
	}

	if bandera && (palabrasRenglon[0] == "*" || palabrasRenglon[0] == "/") {
		llenarLog("ERROR:\tError en la sintaxis de comentario, en la línea: " + strconv.Itoa(num))
		resp = false
	} else {
		if !bandera && (palabrasRenglon[0] == "*" || palabrasRenglon[0] == "/") {
			llenarLog("ERROR:\tError en la sintaxis de comentario, en la línea: " + strconv.Itoa(num))
			resp = true
		}
	}
	return resp
}

/*
# def: Función que sirve para separar código de comentarios que están en la misma línea
# return: Devuelve el código que se extrajo de la línea
# data (in): Parametro que contendrá la línea que se necesita separar
# index: Contiene el número de caracter donde se identificó el símbolo
*/
func extraerCodigo(data string) string {
	var codigo string
	index := strings.IndexAny(data, "-/")
	letras := Tokenizador([]byte(data), "")
	if index != -1 {
		for i := 0; i < len(letras); i++ {
			if i < index-1 {
				if codigo != "" {
					codigo += letras[i]
				} else {
					codigo = letras[i]
				}
			}
		}
	}
	return codigo
}

/*
# def: Función para separar la data según un delimitador que especificamos
# return: Devuelve la data separada en un arreglo de strings
# data (in): Información que le queremos aplicar el proceso descrito
# delimitador (in): Caracter o caracteres por los que queremos separar la data
*/
func Tokenizador(data []byte, delimitador string) []string {
	return strings.Split(string(data), delimitador)
}

/*
# def: Función que sirve para generar los archivos necesarios
# data (in): Parametro que contiene la información que contendrá el archivo
# nombreFile (in): Parametro que contendrá el nombre del archivo
*/
func generarArchivos(data []byte, nombreFile string) {
	path := "../data_sources/" + nombreFile
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		llenarLog("ERROR:\tEl archivo: " + nombreFile + " no se pudo crear")
	} else {
		llenarLog("LOG:\tEl archivo: " + nombreFile + " fue guardado en la ruta: " + path)
	}
}
