package main

import (
    "rdvc/init_dir"
    "rdvc/networking"
	"os"
    "strings"

	"github.com/pterm/pterm"
)

func printHelp() {
    pterm.FgGreen.Println("rdvc - RainDrops Version Control")
    pterm.FgBlue.Println("Usage:")
    pterm.FgCyan.Println(" rdvc init -p <path_to_directory> -n <repo_name>    Initialize a controlled directory")
    pterm.FgCyan.Printfln(" rdvc keep -m <message> -u <user_name> -n <repo_name>    Keep changes with a message")
    pterm.FgCyan.Printfln(" rdvc line -n <repo_name> -o  <name_line>     Create line for current changes in repo")
    pterm.FgCyan.Printfln(" rdvc checkout -n <repo_name> -o <line_name>      Go to line")
    pterm.FgCyan.Printfln(" ")
    pterm.FgCyan.Printfln(" rdvc set -u <username> -p <password>     Setup config for cloud storage")
    pterm.FgCyan.Printfln(" rdvc send -n <repo_name>     Send keeped changes to cloud")
    pterm.FgCyan.Printfln(" rdvc get -n <repo_name>     Get list changes from repo in cloud")
    pterm.FgCyan.Printfln(" rdvc roll -n <repo_name> -i <number_modification>     Pullback to choosed version of file")
    pterm.FgCyan.Printfln(" ")
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
    case "line":
        repoPath := init_dir.ReadFromReg(os.Args[3])
        vcs := init_dir.NewVCS(repoPath)
        err := vcs.CreateBranch(os.Args[5])
        if err != nil {
            pterm.Error.Println("Error creating line:", err)
            return
        }
    case "checkout":
        repoPath := init_dir.ReadFromReg(os.Args[3])
        vcs := init_dir.NewVCS(repoPath)
        err := vcs.CheckoutBranch(os.Args[5])
        if err != nil {
            pterm.Error.Println("Ошибка переключения на ветку:", err)
            return
        }
    case "send":
        networking.UploadKeeps(os.Args[3])
    case "get":
        networking.GetKeeps(os.Args[3])
    case "set":
        networking.Connect(os.Args[3], os.Args[5])

        reg_data := []string{os.Args[3], os.Args[5]}
        init_dir.CreateSettings( strings.Join(reg_data, ","), "reg_data" )
    case "roll":
        
    default:
        pterm.Error.Printfln("Unknown command: %s\n", os.Args[1])
        printHelp()
    }

}
func main() {
    check_args()
}