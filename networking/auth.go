package networking
import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net"
	"regexp"

	"github.com/pterm/pterm"
)
var IsLetter = regexp.MustCompile(`1`).MatchString

var PASSWD string
var UNAME string

func getUserCert(conn net.Conn, username string) bool {
	netbuff := make([]byte, 1024)
	conn.Write([]byte(username + "\n"))

	n, err := conn.Read(netbuff)
	if err != nil {

		return false
	}
	if string(netbuff[:n]) == "1" {

		return true
	} else {
		//pterm.FgRed.Println("Certificate not found")

		return false
	}

}
func AuthenticateClient(conn net.Conn, pass string, uname string) error {

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		pterm.Warning.Println("Не удалось связаться с системой\nИзмените конфигурационный файл или попробуйте снова")
		return err
	}
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgLightGreen)).Println(string(buffer[:n]))

	UNAME = uname
	if getUserCert(conn, uname) {
		//pterm.Success.Println("Cert found")
		return nil
	} else {

		hash := md5.Sum([]byte(pass))
		strPasswd := hex.EncodeToString(hash[:])
		conn.Write([]byte(strPasswd + "\n"))

		n, err = conn.Read(buffer)
		if err != nil {
			return err
		}

		if IsLetter(string(buffer[:n])) {
			PASSWD = pass
			if len(PASSWD) == 0 {
				pterm.FgRed.Println("Error cretating crypto key")
				hash := md5.Sum([]byte(uname))
				strPasswd := hex.EncodeToString(hash[:])

				PASSWD = strPasswd
				pterm.FgBlue.Println("Generated crypto key " + PASSWD)
			}

			return nil
		} else {
			pterm.FgRed.Println("No valid login/password ")
			return errors.New("no valid login/password ")
		}
	}
}
