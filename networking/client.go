package networking

import (
	"crypto/tls"
	"flag"
	"net"
	"strings"

	"github.com/pterm/pterm"

)

const (
	PORT = "2121"
	HOST = "localhost"
)
func SetupRepoInServer(conn net.Conn){
	conn.Write([]byte("set_repo\n"))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		pterm.Error.Println("Error in creating storage: ", err)
		return
	}

	str := strings.Trim(string(buf[:n]), "\n")
	commandArr := strings.Fields(str)
	if commandArr[0] != "Success" {
		pterm.Error.Println("Network error: ", commandArr[1])
		return
	}
	pterm.Success.Println("Repository storage created successfully")

}
func Connect(uname string, passwd string) (err error) {

	var connect net.Conn

	boolTSL := flag.Bool("tls", true, "Set tls connection")
	flag.Parse()
	if !*boolTSL {

		connect, err = net.Dial("tcp", HOST+":"+PORT)
		if err != nil {
			pterm.Warning.Println("Failed to connect to server")
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

	defer connect.Close()

	if err := AuthenticateClient(connect, passwd, uname); err != nil {
		return err
	}
	SetupRepoInServer(connect)
	connect.Write([]byte("close\n"))
	return nil

}
