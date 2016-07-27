package networking

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func init() {
}

// StartServer starts listening for incoming connections
func StartServer() {
	fmt.Println("Starting HTTP Server...")
	http.HandleFunc("/", hello)
	http.HandleFunc("/test", test)
	http.ListenAndServe(":8000", nil)
}

func StartTCPServer() {
	fmt.Println("Starting TCP Server...")
	ln, _ := net.Listen("tcp", ":4242")
	conn, _ := ln.Accept()
	for {
		fmt.Println("Incoming connection:", conn.RemoteAddr())
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
		if strings.HasPrefix(message, "exit") {
			fmt.Println("Closing connection...")
			conn.Close()
			break
		}
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HANDLER Test")
	io.WriteString(w, "Test handler!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("HANDLER Hello")
	fmt.Println(r.Method, " ", r.Host, " ")
	io.WriteString(w, "Hello world!")
}
