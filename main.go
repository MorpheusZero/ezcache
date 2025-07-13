package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"strings"
)

var data = make(map[string]string)

func encode(input string) string {
	// Encode the input string to Base64
	encoded := base64.StdEncoding.EncodeToString([]byte(input))

	fmt.Println("Encoded message:", encoded)

	return encoded
}

func decode(encodedString string) string {
	// Decode the Base64 string
	decoded, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		log.Fatal("Base64 decoding failed:", err)
	}

	decodedStr := string(decoded)

	fmt.Println("Decoded message:", decodedStr)

	return decodedStr
}

func main() {

	// encode("dXNlcmlkXzE6ZHlsYW5sZWdlbmRyZTA5QGdtYWlsLmNvbQ==")

	// Start listening on TCP port 7789
	listener, err := net.Listen("tcp", ":7789")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 7789...")

	for {
		// Accept an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Printf("New connection from %s\n", conn.RemoteAddr())

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("Received: %s\n", text)
		var decodedStr string = decode(text)
		fmt.Printf("Decoded string: %s\n", decodedStr)

		parts := strings.SplitN(decodedStr, ":", 2)
		n := len(parts)
		if n == 2 {
			key := parts[0]
			value := parts[1]
			data[key] = value
			fmt.Printf("Stored key: %s, value: %s\n", key, value)
			conn.Write([]byte("+OK\r\n"))
		} else {
			fmt.Println("Invalid format, expected key:value")
			conn.Write([]byte("-ERR Invalid format\r\n"))
			continue
		}

		if len(data) > 0 {
			for k, v := range data {
				fmt.Printf("Key: %s, Value: %s\n", k, v)
			}
		} else {
			fmt.Println("No data stored yet.")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Connection error:", err)
	}
}
