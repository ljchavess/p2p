package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

var ipServer = ""  /*Variavel para armazenar o IP do servidor*/
var porta = "8080" /*Porta que será armazenado o servidor*/

/*Definição da estrutura que armazenará os usuários conectados ao servidor*/
type Client struct {
	End       string   `json:"end"`
	ListaArqs []string `json:"listaArqs"`
}

var listaClients []Client /*Essa variavel é para armazenar a lista de clientes conectadas ao servidor*/

func handler(conn net.Conn) {
	/*Antes do loop infinito, faz a primeira leitura das infos do cliente */
	m, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		print(err)
	}
	fmt.Println(m)

	registro(conn, m)
	fmt.Println(listaClients)

	/*Laço infinito*/
	for {
		m, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%v Conexão Fechada\n", conn.RemoteAddr())
				conn.Close()
				return
			}
			fmt.Println("Erro de Leitura", err)
			return
		}
		_, err = conn.Write([]byte(m))
		if err != nil {
			fmt.Println("Error writing to connection")
			return
		}
		fmt.Printf("%v %q\n", conn.RemoteAddr(), m)
	}
}

/*Função para executar alguns comandos no servidor*/
func commands() {
	var comando string
	for {
		fmt.Scanf("%s", comando)
		if comando == "/allClients\n" {
			print(listaClients)
		}
	}
}

/*
Função para pegar o IP local do servidor
Será usado como informação para os clientes se conectarem
*/
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

/*É aqui que o programa começa, inicializa algumas configurações do servidor*/
func init() {
	ipServer = GetLocalIP()
}

/*Função de registro do usuário no servidor*/
func registro(conn net.Conn, m string) {
	var client Client
	client.End = conn.RemoteAddr().String()
	m = strings.Replace(m, "\n", "", 1)
	client.ListaArqs = strings.Split(m, ",")
	listaClients = append(listaClients, client)
}

/*Inicio do programa de fato*/
func main() {
	var config = ipServer + ":" + porta

	ln, _ := net.Listen("tcp", config) /*abertuda do servidor da porta 8080*/
	go commands()                      /*abre uma goroutine para executar comandos de debug no servidor*/
	fmt.Println("Servidor escutando no endereço => ", ln.Addr().String())

	for {
		conn, _ := ln.Accept()
		fmt.Println("Connection accepted")
		go handler(conn)
	}
}
