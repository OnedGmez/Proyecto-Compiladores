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
	"strings"
)

// Definición de HashTable

/*
#tamHashTable: Representa el tamaño del arreglo utilizado como hashtable, el tamaño es representado por la cantidad de palabras reservadas en SQL * 2
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

func main() {
	llenarTS()
	fmt.Println(tablaSimbolos.searchInHashTable("INT"))
}

//Funciones para analisis

/*
# def: Función que nos ayudará a cargar la tabla de símbolos en la primera carga del analizador
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

func Tokenizador(data []byte, delimitador string) []string {
	return strings.Split(string(data), delimitador)
}
