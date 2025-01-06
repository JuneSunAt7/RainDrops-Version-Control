package networking

import (
	"crypto/tls"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"rdvc/init_dir"
	"strings"

	"github.com/pterm/pterm"
)
func GetLastKeepFromCloud(repoName string) error {
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
	defer connect.Close()

	connect.Write([]byte(fmt.Sprintf("get_last_keep %s\n", repoName)))
    _, err = connect.Write([]byte("200\n"))
    if err != nil {
        pterm.Error.Printfln("Error sending command to server: %v", err)
    }
	fileName, err := getFileFromServer(connect, repoName)
    if err != nil {
        pterm.Error.Printfln("Error getting file from server: %v", err)
    }

    pterm.Success.Printf("Received file: %s\n", fileName)

	connect.Write([]byte("close\n"))
	return nil
}

var validFileName = map[rune]struct{}{
    '_': {}, 'a': {}, 'b': {}, 'c': {}, 'd': {}, 'e': {}, 'f': {},
    'g': {}, 'h': {}, 'i': {}, 'j': {}, 'k': {}, 'l': {}, 'm': {},
    'n': {}, 'o': {}, 'p': {}, 'q': {}, 'r': {}, 's': {}, 't': {},
    'u': {}, 'v': {}, 'w': {}, 'x': {}, 'y': {}, 'z': {},
    'A': {}, 'B': {}, 'C': {}, 'D': {}, 'E': {}, 'F': {},
    'G': {}, 'H': {}, 'I': {}, 'J': {}, 'K': {}, 'L': {},
    'M': {}, 'N': {}, 'O': {}, 'P': {}, 'Q': {}, 'R': {},
    'S': {}, 'T': {}, 'U': {}, 'V': {}, 'W': {}, 'X': {},
    'Y': {}, 'Z': {}, '0': {}, '1': {}, '2': {}, '3': {},
    '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {},
    '.': {}, '-': {}, 
}

func getValidFileName(name string) string {
    var sb strings.Builder
    for _, r := range name {
        if _, ok := validFileName[r]; ok {
            sb.WriteRune(r)
        } else {
            sb.WriteRune('_')
        }
    }
    return sb.String()
}

func getFileFromServer(conn net.Conn, repoName string) (string, error) {
	dir := init_dir.ReadFromReg(repoName) + "\\.rdvc\\keeps\\"
  
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        return "", fmt.Errorf("error reading file name from server: %v", err)
    }


    fileName := string(buf[:n])
    
    fileName = fileName[:len(fileName)-2]

    
    fileName = getValidFileName(fileName)
	fileName = strings.ReplaceAll(fileName, " ", "_")
    
    var fileSize int64
    err = binary.Read(conn, binary.BigEndian, &fileSize)
    if err != nil {
        return "", fmt.Errorf("error reading file size from server: %v", err)
    }

    
    file, err := os.Create(dir + fileName)
    if err != nil {
        return "", fmt.Errorf("error creating file on client: %v", err)
    }
    defer file.Close()

    
    receivedBytes := int64(0)
    for receivedBytes < fileSize {
        bytesToRead := fileSize - receivedBytes
        if bytesToRead > int64(len(buf)) {
            bytesToRead = int64(len(buf))
        }

        n, err := conn.Read(buf[:bytesToRead])
        if err != nil {
            return "", fmt.Errorf("error reading file data from server: %v", err)
        }

        
        _, err = file.Write(buf[:n])
        if err != nil {
            return "", fmt.Errorf("error writing to file: %v", err)
        }

        receivedBytes += int64(n)
    }

    pterm.Success.Printf("File %s successfully received. Total size: %d bytes.\n", fileName, receivedBytes)
    return fileName, nil
}