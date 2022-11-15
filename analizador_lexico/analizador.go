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
var nuevoElemento [1][5]string
var log string

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

func llenarLog(data string) {
	if log != "" {
		log += "\n" + data
	} else {
		log = data
	}
}

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
		llenarLog("ERROR:\tData vacía, verifique sus archivos y vuelva a intentar")
	}
	return retorno
}

/*
# def: Función que nos ayudará a cargar la tabla de símbolos en la primera carga del analizador
# dataTS: array que nos sirve para acomodar la información de palabras propias de SQL que serán guardadas en la tabla de simbolos
# nombreFile: Tiene seteado el nombre del archivo que contiene las palabras propias de SQL
# path: Tiene seteada la ruta del archivo ue contiene las palabras propias de SQL
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

func AnalizarSQL(nombreSQL string, Op int) string {
	numLinea = 0
	var identificador string
	var delim string
	var nuevaDataTS [1][5]string
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
				reIdentificador := "^@.*$"
				reReservada := "^[^@]\\D*$"
				reIdentificadorND := "^[^@].*$"
				palabrasLinea := Tokenizador([]byte(renglonSQL), " ")
				for _, palabra := range palabrasLinea {
					if strings.TrimSpace(palabra) != "" {
						letras := Tokenizador([]byte(palabra), "")
						for _, letra := range letras {
							runaLetra = rune(letra[0])

							if unicode.IsLetter(runaLetra) || (unicode.IsNumber(runaLetra) && identificador != "") || runaLetra == '_' || runaLetra == '@' {

								if identificador != "" {
									identificador += string(runaLetra)
								} else {
									identificador = string(runaLetra)
								}
							}

							if (unicode.IsGraphic(runaLetra) && runaLetra != '_' && runaLetra != '@' && runaLetra != '.') || runaLetra == '*' || runaLetra == ',' {
								if delim != "" {
									delim += string(runaLetra)
								} else {
									guardarDelimitador(string(runaLetra))
									delim = ""
								}
							}

						}

						if identificador != "" {
							coincide, _ := regexp.Match(reReservada, []byte(identificador))
							if coincide {
								//fmt.Println("Reservada: " + identificador)
								if tablaSimbolos.searchInHashTable(strings.ToUpper(identificador)) {
									tablaSimbolos.modificarInHashTable(strings.ToUpper(identificador), strconv.Itoa(numLinea+1))
									tipo := tablaSimbolos.searchTipoInHashTable(strings.ToUpper(identificador))
									insertarTokenLexema(strings.ToUpper(identificador), tipo, numLinea+1)
								} else {
									insertarTokenLexema(identificador, "Identificador", numLinea+1)
									nuevaDataTS[0][0] = identificador
									nuevaDataTS[0][1] = "Identificador"
									nuevaDataTS[0][2] = strconv.Itoa(numLinea + 1)
									nuevaDataTS[0][3] = "NE"
									nuevaDataTS[0][4] = "Usuario"
									tablaSimbolos.insertInHashTable(nuevaDataTS)
								}
							} else {
								coincide, _ := regexp.Match(reIdentificadorND, []byte(identificador))
								if coincide {
									//fmt.Println("Identificador: " + identificador)
									if tablaSimbolos.searchInHashTable(identificador) {
										tablaSimbolos.modificarInHashTable(identificador, strconv.Itoa(numLinea+1))
										tipo := tablaSimbolos.searchTipoInHashTable(identificador)
										insertarTokenLexema(identificador, tipo, numLinea+1)
									} else {
										insertarTokenLexema(identificador, "Identificador", numLinea+1)
										nuevaDataTS[0][0] = identificador
										nuevaDataTS[0][1] = "Identificador"
										nuevaDataTS[0][2] = strconv.Itoa(numLinea + 1)
										nuevaDataTS[0][3] = "NE"
										nuevaDataTS[0][4] = "Usuario"
										tablaSimbolos.insertInHashTable(nuevaDataTS)
									}
								} else {
									coincide, _ := regexp.Match(reIdentificador, []byte(identificador))
									if coincide {
										//fmt.Println("Declaración: " + identificador)
										if tablaSimbolos.searchInHashTable(identificador) {
											tablaSimbolos.modificarInHashTable(identificador, strconv.Itoa(numLinea+1))
											tipo := tablaSimbolos.searchTipoInHashTable(identificador)
											insertarTokenLexema(identificador, tipo, numLinea+1)
										} else {
											insertarTokenLexema(identificador, "Declaración de variable", numLinea+1)
											nuevaDataTS[0][0] = identificador
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
						identificador = ""
					}

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

func guardarDelimitador(delim string) {
	reSimbolo := "\\W"
	coincide, _ := regexp.Match(reSimbolo, []byte(delim))
	if coincide {
		//fmt.Println("Declaración: " + identificador)
		if tablaSimbolos.searchInHashTable(delim) {
			tablaSimbolos.modificarInHashTable(delim, strconv.Itoa(numLinea+1))
			tipo := tablaSimbolos.searchTipoInHashTable(delim)
			insertarTokenLexema(delim, tipo, numLinea+1)
		} else {
			llenarLog("WARNING:\tSimbolo: " + delim + ", en la línea: " + strconv.Itoa(numLinea+1) + ", no identificado dentro de la sintaxís definida para SQL")
			insertarTokenLexema("[W]"+delim, "Delimitador/Separador", numLinea+1)
		}
	}
}

func insertarTokenLexema(data string, tipo string, num int) string {
	if tokensLexemas != "" {
		tokensLexemas += "\n" + data + ", " + tipo + ", " + strconv.Itoa(num)
	} else {
		tokensLexemas = data + ", " + tipo + ", " + strconv.Itoa(num)
	}
	return tokensLexemas
}

func ExisteComentario(data string) bool {
	return strings.ContainsAny(data, "-/*")
}

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

func llenarPosibleComentario(data string) {
	if data != "" {
		if posibleComentario != "" {
			posibleComentario += "\n" + data
		} else {
			posibleComentario = data
		}
	}
}

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

func generarArchivos(data []byte, nombreFile string) {
	path := "../data_sources/" + nombreFile
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		llenarLog("ERROR:\tEl archivo: " + nombreFile + " no se pudo crear")
	} else {
		llenarLog("LOG:\tEl archivo: " + nombreFile + " fue guardado en la ruta: " + path)
	}
}
