package main

import (
        "io"
        "log"
        "net"
        "sync"
        "github.com/buseliiik/is105sem03/mycrypt"
)

func main() {
        var wg sync.WaitGroup

        server, err := net.Listen("tcp", "172.17.0.3:8002")
        if err != nil {
                log.Fatal(err)
        }
        log.Printf("bundet til %s", server.Addr().String())
        wg.Add(1)
        go func() {
                defer wg.Done()
                for {
                        log.Println("f  r server.Accept() kallet")
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
                                                return
                                        }
                                        log.Println("Krypter melding: ", string(buf[:n]))
                                        cryptedMsg := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
                                        log.Println("Kryptert melding: ", string(cryptedMsg))

                                        switch msg := string(cryptedMsg); msg {
                                        case "ping":
                                                _, err = c.Write([]byte("pong"))
                                        default:
                                                _, err = c.Write([]byte(string(cryptedMsg)))
                                        }
                                        if err != nil {
                                                if err != io.EOF {
                                                        log.Println(err)                                                }
                                                return
                                        }
                                }
                        }(conn)
                }
        }()
        wg.Wait()
}
