package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var ipServer = ""
var porta = ""
var conn net.Conn

var listaArquivos []string

var (
	output    = make(chan string)
	input     = make(chan string)
	errorChan = make(chan error)
)

func readStdin() {
	for {
		reader := bufio.NewReader(os.Stdin)
		m, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input <- m
	}
}

func readConn(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		m, err := reader.ReadString('\n')
		if err != nil {
			errorChan <- err
			return
		}
		output <- m
	}
}

/*Função persistente de conexão ao servidor */
func connect(configServer string) net.Conn {
	var (
		//conn net.Conn
		err error
	)
	for {
		fmt.Println("Conectando ao servidor...")
		conn, err = net.Dial("tcp", configServer)
		if err == nil {
			break
		}
		fmt.Println(err)
		time.Sleep(time.Second * 1)
	}
	fmt.Println("Conexão Aceita")
	return conn
}

func init() {
}

func main() {
	files, err := os.ReadDir("./files/") /*Leitura dos arquivos que estarão disponiveis*/
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(files)

	for _, f := range files {
		listaArquivos = append(listaArquivos, f.Name())
	}
	//fmt.Println(listaArquivos)
	teste := strings.Join(listaArquivos, ",")
	println(teste)

	fmt.Print("Digite o IP do servidor: ")
	fmt.Scanf("%s", &ipServer)
	fmt.Print("Digite a porta do servidor: ")
	fmt.Scanf("%s", &porta)

	var configServer = ipServer + ":" + porta

	go readStdin()

	conn := connect(configServer)
	_, err = conn.Write([]byte(teste + "\n"))
	if err != nil {
		fmt.Println(err)
	}

RECONNECT:
	for {
		go readConn(conn)

		for {
			select {
			case m := <-output:
				fmt.Printf("Recebido: %q\n", m)

			case m := <-input:
				fmt.Printf("Enviado: %q\n", m)
				_, err := conn.Write([]byte(m + "\n"))
				if err != nil {
					fmt.Println(err)
					conn.Close()
					continue RECONNECT
				}
			case err := <-errorChan:
				fmt.Println("Error:", err)
				conn.Close()
				continue RECONNECT
			}
		}
	}
}
