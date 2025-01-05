package networking

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"strings"

	"golang.org/x/sys/windows/registry"

	"rdvc/init_dir"

	"github.com/pterm/pterm"

	"encoding/binary"
	"io"
	"os"
	"path/filepath"
)
func getCreds() []string{
	key := registry.CURRENT_USER
	subKey := "Software\\RaindropsVC"
	valueName := "reg_data"

	val, err := init_dir.ReadRegistryValue(key, subKey, valueName)
		if err != nil {
			pterm.Error.Println("Error connect:", err)
			return []string{}
		}
	return strings.Split(val, ",")
	
}
func UploadKeeps(repoName string) error {
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
	fmt.Println(reg_data)

	defer connect.Close()


	if err := AuthenticateClient(connect, reg_data[1], reg_data[0]); err != nil {
		return err
	}
	uploadFiles(connect, repoName)
	connect.Write([]byte("close\n"))
	return nil
}
func uploadFiles(conn net.Conn,  repoName string) {
	sourceFolder := init_dir.ReadFromReg(repoName) + "\\.rdvc\\keeps"  
    
    conn.Write([]byte(fmt.Sprintf("update_repo %s\n\x00",repoName))) 

    files, err := os.ReadDir(sourceFolder)
    if err != nil {
        pterm.Error.Printfln("Error reading dir: %v", err)
    }
	

    for _, file := range files {
        if file.IsDir() {
            continue 
        }

        filePath := filepath.Join(sourceFolder, file.Name())
        
        
 		fileInfo, err := os.Stat(filePath)
        if err != nil {
            pterm.Error.Printfln("Error getting info about file %s: %v", file.Name(), err)
            continue
        }
		buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		pterm.Error.Println("Error uploading file: ", err)
		return
	}

	str := strings.Trim(string(buf[:n]), "\n")
	commandArr := strings.Fields(str)
	if commandArr[0] != "200" {
		pterm.Error.Println("Error network connection")
	}

        _, err = conn.Write([]byte(fmt.Sprintf("%s\x00\n", file.Name()))) // Нулевой байт в конце имени
        if err != nil {
            pterm.Error.Printfln("Error getting info about file: %v", err)
            continue
        }

        size := fileInfo.Size()
        err = binary.Write(conn, binary.BigEndian, size)
        if err != nil {
             pterm.Error.Printfln("Error getting info about file: %v", err)
            continue
        }

        f, err := os.Open(filePath)
        if err != nil {
             pterm.Error.Printfln("Error opening file %s: %v", file.Name(), err)
            continue
        }
        defer f.Close()


        _, err = io.Copy(conn, f)
        if err != nil {
             pterm.Error.Printfln("Error upload file %s: %v", file.Name(), err)
            continue
        }
		
		pterm.Success.Printfln("File %s succesfuly clouding.", file.Name())
    }


    _, err = conn.Write([]byte("\x00")) 
    if err != nil {
         pterm.Error.Printfln("Error: %v", err)
    }
}