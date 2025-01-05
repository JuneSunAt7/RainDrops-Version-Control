package main

import (
    "rdvc/init_dir"
    "rdvc/networking"
    "os"
    "strings"

    "github.com/pterm/pterm"
)

var (
    currentRepoName string
    currentRepoPath string
)

func printHelp() {
    pterm.FgGreen.Println("rdvc - RainDrops Version Control")
    pterm.FgBlue.Println("Usage:")
    pterm.FgCyan.Println(" rdvc init -p <path_to_directory> -n <repo_name>    Initialize a controlled directory")
    pterm.FgCyan.Printfln(" rdvc keep -m <message> -u <user_name>    Keep changes with a message")
    pterm.FgCyan.Printfln(" rdvc line -o <name_line>      Create line for current changes in repo")
    pterm.FgCyan.Printfln(" rdvc checkout -o <line_name>      Go to line")
    pterm.FgCyan.Printfln(" rdvc roll      Pullback to choosed version of file")
    pterm.FgCyan.Printfln(" ")
    pterm.FgCyan.Printfln(" rdvc session -u <user_name> -n <repo_name>     Start work session and store data")
    pterm.FgCyan.Printfln(" ")
    pterm.FgCyan.Printfln(" rdvc set -u <username> -p <password>     Setup config for cloud storage")
    pterm.FgCyan.Printfln(" rdvc send     Send kept changes to cloud")
    pterm.FgCyan.Printfln(" rdvc get     Get list changes from repo in cloud")
   
    pterm.FgCyan.Printfln(" ")
    pterm.FgMagenta.Printfln(" rdvc help     Display this help message")
}

func checkArgs() {
    if len(os.Args) < 2 {
        printHelp()
        return
    }
    switch os.Args[1] {
    case "help":
        printHelp()
        return
    case "init":
        pathFlag := os.Args[3]
        init_dir.InitInvisible(pathFlag)
        init_dir.CreateSettings(pathFlag, os.Args[5])
    case "session":
        if len(os.Args) < 5 {
            pterm.Error.Println("You must specify a user name and repo name for the session.")
            return
        }
        currentRepoName = os.Args[5]
        currentRepoPath = init_dir.ReadFromReg(currentRepoName)
        // Вы здесь можете также запустить какой-либо код для начала сессии

    case "keep":
        if currentRepoPath == "" {
            pterm.Error.Println("You must start a session before using this command.")
            return
        }
        versionControl := init_dir.NewVCS(currentRepoPath)
        message := os.Args[3]
        author := os.Args[5]

        err := versionControl.MakeKeep(message, author)
        if err != nil {
            pterm.Error.Println("Error in keeping:", err)
        }
    case "line":
        if currentRepoPath == "" {
            pterm.Error.Println("You must start a session before using this command.")
            return
        }
        vcs := init_dir.NewVCS(currentRepoPath)
        err := vcs.CreateBranch(os.Args[3])
        if err != nil {
            pterm.Error.Println("Error creating line:", err)
            return
        }
    case "checkout":
        if currentRepoPath == "" {
            pterm.Error.Println("You must start a session before using this command.")
            return
        }
        vcs := init_dir.NewVCS(currentRepoPath)
        err := vcs.CheckoutBranch(os.Args[3])
        if err != nil {
            pterm.Error.Println("Error switching branch:", err)
            return
        }
    case "send":
        // Если repoName не требуется, используем текущий
        if currentRepoName == "" {
            pterm.Error.Println("You must start a session to send changes.")
            return
        }
        networking.UploadKeeps(currentRepoName)
    case "get":
        // Если repoName не требуется, используем текущий
        if currentRepoName == "" {
            pterm.Error.Println("You must start a session to get changes.")
            return
        }
        networking.GetKeeps(currentRepoName)
    case "set":
        networking.Connect(os.Args[3], os.Args[5])
        reg_data := []string{os.Args[3], os.Args[5]}
        init_dir.CreateSettings(strings.Join(reg_data, ","), "reg_data")
    case "roll":
        if currentRepoPath == "" {
            pterm.Error.Println("You must start a session before using this command.")
            return
        }
        vcs := init_dir.NewVCS(currentRepoPath)
        err := vcs.Rollback()
        if err != nil {
            pterm.Error.Println("Error pullback:", err)
            return
        }
    default:
        pterm.Error.Printfln("Unknown command: %s\n", os.Args[1])
        printHelp()
    }
}

func main() {
    checkArgs()
}