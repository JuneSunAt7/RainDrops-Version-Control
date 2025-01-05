package networking

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"

	"github.com/pterm/pterm"

)
func GetLastKeep(repoName string) error {
	var connect net.Conn
	var err error

	boolTSL := flag.Bool("tls", true, "Set tls connection")
	flag.Parse()
	if !*boolTSL {
		connect, err = net.Dial("tcp", HOST+":"+PORT)
		if err != nil {
			return err
		}
	} else {
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}

		connect, err = tls.Dial("tcp", HOST+":"+PORT, conf)
		if err != nil {
			return err
		}
	}
	reg_data := getCreds()

	defer connect.Close()


	if err := AuthenticateClient(connect, reg_data[1], reg_data[0]); err != nil {
		return err
	}
	gettingLastKeep(connect, repoName)
	connect.Write([]byte("close\n"))
	return nil
}
func gettingLastKeep(conn net.Conn, repoName string){
	conn.Write([]byte(fmt.Sprintf("get_last_keep %s\n", repoName)))
	buffer := make([]byte, 4096)
	n, _ := conn.Read(buffer)
	pterm.FgGreen.Println(string(buffer[:n]))
}