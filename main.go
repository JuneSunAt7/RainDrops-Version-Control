package main

import (
    "rdvc/init_dir"
    "rdvc/networking"
	"os"

	"github.com/pterm/pterm"
)

func printHelp() {
    pterm.FgGreen.Println("rdvc - RainDrops Version Control")
    pterm.FgBlue.Println("Usage:")
    pterm.FgCyan.Println(" rdvc init -p <path_to_directory> -n <nickname_repo>    Initialize a controlled directory")
    pterm.FgCyan.Printfln(" rdvc keep -m <message> -u <user_name> -n <nickname_repo>    Keep changes with a message")
    pterm.FgCyan.Printfln(" rdvc send -n <nickname_repo>     Send keeped changes to cloud")
    pterm.FgCyan.Printfln(" rdvc get -n <name_dir_in_cloud>     Get last changes from cloud")
    pterm.FgCyan.Printfln(" ")
    pterm.FgCyan.Printfln(" rdvc set -u <username> -p <password>     Setup config for cloud storage")

    pterm.FgMagenta.Printfln(" rdvc help     Display this help message")
}
func check_args(){
    if len(os.Args) < 2 {
        printHelp()
        return
    }
    switch os.Args[1]{
    case "help":
        printHelp()
        return
    case "init":
        pathFlag := os.Args[3]
        init_dir.InitInvisible(pathFlag)
        init_dir.CreateSettings(pathFlag, os.Args[5])
    case "keep":
        repoPath := init_dir.ReadFromReg(os.Args[7])
        versionControl := init_dir.NewVCS(repoPath)
        message := os.Args[3]
        author := os.Args[5]

        err := versionControl.MakeKeep(message, author)
        if err != nil {
            pterm.Error.Println("Error in keeping:", err)
        }
    case "send":
        
    case "get":
    case "set":
        networking.Connect(os.Args[3], os.Args[5])
    default:
        pterm.Error.Printfln("Unknown command: %s\n", os.Args[1])
        printHelp()
    }
}
func main() {
    check_args()
}