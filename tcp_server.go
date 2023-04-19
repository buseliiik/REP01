package main

import (
	"io"
	"log"
	"net"
	"sync"
	"github.com/buseliiik/is105sem03/mycrypt"
	 "github.com/buseliiik/funtemps/conv" 
	"github.com/buseliiik/minyr/yr"
)

func main() {
	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for-løkke
					}
					log.Println("Krypter melding: ", string(buf[:n]))
					cryptedMsg := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Kryptert melding: ", string(cryptedMsg))
					
					switch msg := string(cryptedMsg); msg {
case "ping":
    _, err = c.Write([]byte("pong"))
case "Kjevik":
    temp, err := yr.GetTemperature()
    if err != nil {
        log.Println(err)
        return
    }
    fahrTemp := conv.CelsiusToFahrenheit(temp)
    resp := fmt.Sprintf("Temperature in Kjevik: %.2f F", fahrTemp)
    _, err = c.Write([]byte(resp))
default:
    _, err = c.Write(cryptedMsg)
}

					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for-løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
