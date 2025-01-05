package init_dir

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pterm/pterm"
)

type Keep struct {
    Message string    `json:"message"`
    Author  string    `json:"author"`
    Date    time.Time `json:"date"`
    Files   []File    `json:"files"`
}

type File struct {
    Name    string `json:"name"`
    Content string `json:"content"`
}

type VCS struct {
    RepoPath string
}

func NewVCS(repoPath string) *VCS {
    return &VCS{RepoPath: repoPath}
}

func (v *VCS) MakeKeep(message string, author string) error {
    keep := Keep{
        Message: message,
        Author:  author,
        Date:    time.Now(),
    }

    // get files in the repository
    err := filepath.Walk(v.RepoPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        
        content, err := ioutil.ReadFile(path)
        if err != nil {
            return err
        }

        // add to keep
        keep.Files = append(keep.Files, File{
            Name:    path,
            Content: string(content),
        })
        return nil
    })

    if err != nil {
        return err
    }

    // save keep to file
    keepFileName := fmt.Sprintf("%d.keep.json", time.Now().Unix())
    keepFilePath := filepath.Join(v.RepoPath, ".rdvc", "keeps", keepFileName)

    os.MkdirAll(filepath.Dir(keepFilePath), os.ModePerm) // dir for keep file

    jsonData, err := json.Marshal(keep)
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(keepFilePath, jsonData, 0644)
    if err != nil {
        return err
    }

    pterm.Success.Printfln("Success!Saved to:  %s\n", keepFilePath)
    return nil
}