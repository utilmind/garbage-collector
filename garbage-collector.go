/*
Garbage Collector - CLI tool to delete outdated files, expired for more than N days
PLEASE BE CAREFUL! AUTHORS ARE NOT RESPONSIBLE IF YOU ACCIDENTALLY DELETE SOMETHING IMPORTANT!

Version: 1.0.0
Copyright (c) 2024 https://github.com/utilmind/

This software is licensed under the MIT License.
You can find the full license text at https://opensource.org/licenses/MIT

    TODO
        1. Complete regular script.
        2. Add support for multiple directories.
*/
package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "github.com/fatih/color"
    // "time"
)

// @CONFIG
const DefExpireDays = 90 // Default number of days before a file is considered expired. Used if -expire argument is omitted.

var cliArgs = map[string]*string{
    "dir":     flag.String("dir", "", "directory name (required)"),
    "sub":     flag.String("sub", "", "include subdirectories"),
    "ext":     flag.String("ext", "", "file extension(s). Comma-separated if multiple"),
    "expire":  flag.String("expire", strconv.Itoa(DefExpireDays), "expire after N days. 0 = don’t check date, delete all"), // AK: actually it's integer, but I'd prefer to parse it myself
    "confirm": flag.String("confirm", "", "'y' or 'yes' auto-confirms file deletions. Otherwise you’ll need to confirm file deletions one by one"),
    "silent":  flag.String("silent", "", "don’t show the names of deleted files, if deletion is auto-confirmed (by -confirm=yes option)"),
}

// @private functions
func showError() {
    errorColor := color.New(color.FgRed, color.Bold)
    errorColor.Print("ERROR: ")
}

func die(str string) {
    showError()
    fmt.Println(str)
    os.Exit(0)
}

func main() {
    flag.Usage = func() {
        thisExeName := filepath.Base(os.Args[0])

        fmt.Printf("\nUsage of %s:\n", thisExeName)

        // This is instead of flag.PrintDefaults()...
        flag.VisitAll(func(f *flag.Flag) {
            def := f.DefValue
            if "" != def {
                def = " (Default: " + def + ".)"
            }

            fmt.Printf("    -%s\t%s.%s\n", f.Name, f.Usage, def) //, f.Default)
        })

        // Add custom description
        fmt.Printf("\nExample: %s -dir=/var/www/your-project-name/data/cache -ext=jpg,jpeg,png,gif,webp -expire=60\n", thisExeName)
    }

    flag.Parse()

    //workDir := strings.TrimRight(*cliArgs["dir"], "/\\") // strip trailing slashes (works for both types, Linux/Mac and Windows)
    workDir := *cliArgs["dir"]
    if "" == workDir {
        showError()
        fmt.Print("-dir argument is required.\n")
        flag.Usage()
        os.Exit(0)
    }

    if '/' == []rune(workDir)[0] && 2 > strings.Count(workDir, "/") {
        die(fmt.Sprintf("No, it doesn’t works with root or any top-level directory. It will not process \"%s\" or any other \"/directory/\" under root.", workDir))
    }

    expire, err := strconv.Atoi(*cliArgs["expire"])
    if nil != err {
        die(fmt.Sprintf("Invalid integer value in argument -expire=%s. Please use integer value to specify the number of days, or skip it to use default %d days.", *cliArgs["expire"], DefExpireDays))
    }

    // Calculate the expiration time
    expireDuration := time.Duration(expire) * 24 * time.Hour
    expireTime := time.Now().Add(-expireDuration)

    // Walk through the directory and find files older than expireTime
    err = filepath.Walk(workDir, func(path string, info os.FileInfo, err error) error {
        if nil != err {
            return err
        }

        // Check if the file is older than expireTime
        if !info.IsDir() && info.ModTime().Before(expireTime) {
            if "" == *cliArgs["confirm"] {
                reader := bufio.NewReader(os.Stdin)
                fmt.Printf("Do you really want to delete file `%s`? (y/n): ", path)
                response, _ := reader.ReadString('\n')
                response = strings.TrimSpace(strings.ToLower(response))
                if "y" != response && "yes" != response {
                    fmt.Printf("Skipped file `%s`.\n", path)
                    return nil
                }
            }else if "" == *cliArgs["silent"] {
                fmt.Printf("Deleting `%s`...\n", path)
            }

            err := os.Remove(path)
            if nil != err {
                // Don't die, just display error and continue...
                fmt.Printf("Can’t delete file `%s`.\n%v\n", path, err)
            }
        }

        return nil // success
    })

    if nil != err {
        die(fmt.Sprintf("Error walking the path `%s`.\n%v\n", workDir, err))
    }
}
