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
	"strconv"
	"strings"
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
	key  [1][6]string
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

func (h *hashtable) insertInHashTable(key [1][6]string) {
	index := hash(key[0][0])
	h.array[index].insertBucketNode(key)
}

func (b *bucket) insertBucketNode(k [1][6]string) {
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
		fmt.Println("Palabra no encontrada")
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
				nodoModificar.key[0][4] = nodoModificar.key[0][4] + "," + linea
				if nodoAnterior != nil {
					nodoAnterior.next = nodoModificar
				} else {
					nodoAnterior = nodoModificar
				}
				b.head = nodoAnterior
				fmt.Println(nodoModificar)
				break
			}
		}
	} else {
		fmt.Println("No se encontró la Key")
	}
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

func main() {
	var nombreArchivo string
	llenarTS()
	fmt.Println("Analisis iniciado")
	fmt.Println("Ingrese el nombre del archivo con su extensión")
	fmt.Scanln(&nombreArchivo)
	dataLimpia := AnalizarSQL(nombreArchivo, 1)
	nombreArchivo = "dataLimpia_" + nombreArchivo
	generarArchivos([]byte(dataLimpia), nombreArchivo)
	fmt.Println("Abtención de Token, iniciada")
	dataTokensLexemas := AnalizarSQL(nombreArchivo, 2)
	nombreArchivo = strings.Replace(nombreArchivo, "dataLimpia_", "tokensLexemas_", 1)
	nombreArchivo = strings.Replace(nombreArchivo, "sql", "txt", 1)
	generarArchivos([]byte(dataTokensLexemas), nombreArchivo)
}

//Funciones para analisis

/*
# def: Función que nos ayudará a cargar la tabla de símbolos en la primera carga del analizador
# dataTS: array que nos sirve para acomodar la información de palabras propias de SQL que serán guardadas en la tabla de simbolos
# nombreFile: Tiene seteado el nombre del archivo que contiene las palabras propias de SQL
# path: Tiene seteada la ruta del archivo ue contiene las palabras propias de SQL
*/
func llenarTS() {
	var dataTS [1][6]string
	var nombreFile string = "file_TS.csv"
	var path string = "../data_sources/" + nombreFile
	file, err := ioutil.ReadFile(path)
	if err == nil {
		renglonData := Tokenizador([]byte(file), "\n")
		for _, renglon := range renglonData {
			ix := 0
			palabrasRenglon := Tokenizador([]byte(renglon), ",")
			for _, palabra := range palabrasRenglon {
				dataTS[0][ix] = string(palabra)
				ix++
			}
			tablaSimbolos.insertInHashTable(dataTS)
		}
	} else {
		fmt.Println("Archivo " + nombreFile + ", no encontrado")
	}
}

func AnalizarSQL(nombreSQL string, Op int) string {
	numLinea = 0
	var dataRetorno string
	SQL, err := ioutil.ReadFile("../data_sources/" + nombreSQL)
	if err == nil {
		lineasSQL := Tokenizador(SQL, "\n")
		for _, renglonSQL := range lineasSQL {
			switch Op {
			case 1:
				if strings.TrimSpace(renglonSQL) != "" {
					switch ExisteComentario(renglonSQL) {
					case true:
						if !bandera {
							AnalizarComentario(renglonSQL, numLinea)
							dataPrecomentario := extraerCodigo(renglonSQL)
							if dataPrecomentario != "" {
								dataRetorno = string(fmt.Appendln([]byte(dataRetorno), dataPrecomentario))
							}
						} else {
							AnalizarComentario(renglonSQL, numLinea)
						}
						break
					case false:
						if !bandera {
							dataRetorno = string(fmt.Appendln([]byte(dataRetorno), renglonSQL))
						} else {
							posibleComentario = string(fmt.Appendln([]byte(posibleComentario), renglonSQL))
						}
						break
					}
				}
				break
			case 2:
				break
			}
			numLinea++
		}
	} else {
		fmt.Println("El archivo " + nombreSQL + ", No se encontró o no se logró abrir")
	}
	return dataRetorno
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
			fmt.Println("Error en la sintaxis de comentario, en la línea: " + strconv.Itoa(num+1))
			bandera = false
		} else {
			if !aperturaMulti {
				if !cierreMulti {
					if potencialErr2 || potencialErr3 {
						posibleComentario = ""
						if VerificarPotencialError(dataSQL, num) {
							if bandera == true {
								posibleComentario = ""
								bandera = false
							} else {
								bandera = true
							}
						} else {
							if !bandera {
								bandera = true
							} else {
								bandera = false
							}
						}
					}
				} else {
					bandera = false
				}
			} else {
				if cierreMulti {
					posibleComentario = string(fmt.Appendln([]byte(posibleComentario), dataSQL))
					bandera = false
				} else {
					bandera = true
				}
			}
		}
	} else {
		posibleComentario = string(fmt.Appendln([]byte(posibleComentario), dataSQL))
		bandera = false
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
		fmt.Println("Error en la sintaxis de comentario, en la línea: " + strconv.Itoa(num+1))
		resp = false
	} else {
		if !bandera && (palabrasRenglon[0] == "*" || palabrasRenglon[0] == "/") {
			fmt.Println("Error en la sintaxis de comentario, en la línea: " + strconv.Itoa(num+1))
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
	if strings.Contains(nombreFile, "dataLimpia_") {
		err := ioutil.WriteFile("../data_sources/"+nombreFile, data, 0644)
		if err != nil {
			fmt.Println("El archivo " + nombreFile + " no se pudo crear")
		}
	}
}
