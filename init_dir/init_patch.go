package init_dir

import (
    "os"
    "github.com/pterm/pterm"
    "syscall"
    "time"
    "golang.org/x/sys/windows/registry"
)

func InitInvisible(directory string) {
    if directory == "" {
        pterm.Error.Println("Error: No directory specified. Use -p <path_to_directory>.")
        return
    }
    err := os.Mkdir(directory + "\\.rdvc", os.ModePerm) 
    if err != nil {
        pterm.Error.Printfln("Failed to create directory: %v\n", err)
        return
    }

    err = syscall.SetFileAttributes(syscall.StringToUTF16Ptr(directory+"\\.rdvc"), syscall.FILE_ATTRIBUTE_HIDDEN)
    if err != nil {
        pterm.Error.Println("Error setting directory attributes:", err)
        return
    }
    p, _ := pterm.DefaultProgressbar.WithTotal(10).WithTitle("Initializing directory").Start()

    for i := 0; i < p.Total; i++ {
        p.UpdateTitle("Initializing directory...")
        p.Increment()
        time.Sleep(time.Millisecond * 50)
    }
    pterm.Success.Printfln("Successful init: %s\n", directory)
}
func CreateSettings(directory string, nick string) {
    if directory == "" {
        pterm.Error.Println("Error: No directory specified. Use -p <path_to_directory>.")
        return
    }
    CreateSettingsToRegedit(nick, directory)
}
func ReadFromReg(nick string) string {
    value, err := ReadRegistryValue(registry.CURRENT_USER, `Software\RaindropsVC`, nick)
    if err != nil {
        return ""
    }
    return value
}