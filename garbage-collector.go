/*
Garbage Collector - CLI tool to delete outdated files, expired for more than N days
PLEASE BE CAREFUL! AUTHORS ARE NOT RESPONSIBLE IF YOU ACCIDENTALLY DELETE SOMETHING IMPORTANT!

Version: 1.0.0
Copyright (c) 2025 https://github.com/utilmind/

This software is licensed under the MIT License.
You can find the full license text at https://opensource.org/licenses/MIT
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
// ------- end of configuration --------

var (
    cliArgs = map[string]interface{}{
        // -dir can be path to a single file, but this feature is not documented. Because if file already deleted, it will display 'error', that file doesn't exists...
        "dir":     flag.String("dir", "", "directory name (required). For safety reasons, it doesn’t works with top-level directories: the root and the next level below the root."),
        "sub":     flag.Bool("sub", false, "(boolean) include subdirectories, if -sub specified"),
        "ext":     flag.String("ext", "", "file extension(s). Comma-separated if multiple"),
        "expire":  flag.String("expire", strconv.Itoa(DefExpireDays), "expire after N days. 0 = don’t check date, delete all"), // AK: actually it's integer, but I'd prefer to parse it myself
        "confirm": flag.Bool("confirm", false, "(boolean) auto-confirms file deletions. Otherwise you’ll need to confirm file deletions one by one"),
        "silent":  flag.Bool("silent", false, "(boolean) don’t show the names of deleted files, if deletion is auto-confirmed (by -confirm option)"),
    }
    // map doesn't guarantee preserving the items order, so list list all our items in the order we prefer
    flagOrder = []string{"dir", "sub", "ext", "expire", "confirm", "silent"}
)

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

func checkGarbageFile(path string, info os.FileInfo, expireTime time.Time) {
    // Check if the file is older than expireTime
    if info.ModTime().Before(expireTime) {
        if !*cliArgs["confirm"].(*bool) {
            reader := bufio.NewReader(os.Stdin)
            fmt.Printf("Do you really want to delete file `%s`? (y/n): ", path)

            // User input. Wait for 'y' or 'yes', case insensitive...
            response, _ := reader.ReadString('\n')
            response = strings.TrimSpace(strings.ToLower(response))
            if "y" != response && "yes" != response {
                fmt.Printf("Skipped file `%s`.\n", path)
                return
            }
        }else if !*cliArgs["silent"].(*bool) {
            fmt.Printf("Deleting `%s`...\n", path)
        }

        err := os.Remove(path)
        if nil != err {
            // Don't die, just display error and continue...
            fmt.Printf("Can’t delete file `%s`.\n%v\n", path, err)
        }
    }
}

func main() {
    flag.Usage = func() {
        thisExeName := filepath.Base(os.Args[0])

        fmt.Printf("\nUsage of %s:\n", thisExeName)

        // This is instead of flag.PrintDefaults()...
        // ...and preserving original order of arguments, instead of alphabetic; w/o using the `flag.VisitAll(func(f *flag.Flag) { ... }`
        for _, flagName := range flagOrder {
            if f := flag.Lookup(flagName); nil != f {
                def := f.DefValue

                // show default value, but only for string types
                if _, ok := cliArgs[f.Name].(*string); ok && def != "" {
                    def = " (Default: " + def + ".)"
                }else {
                    def = "" // no default for non-string types
                }

                fmt.Printf("    -%s\t%s.%s\n", f.Name, f.Usage, def) //, f.Default)
            }
        }

        // Add custom description
        fmt.Printf("\nExample: %s -dir=/var/www/your-project-name/data/cache -ext=jpg,jpeg,png,gif,webp -expire=60\n", thisExeName)
    }

    flag.Parse()

    workDir := *cliArgs["dir"].(*string) // BTW no need to trim slashes, it works in either case
    if "" == workDir {
        showError()
        fmt.Print("-dir argument is required.\n")
        flag.Usage()
        os.Exit(0)
    }

    if '/' == []rune(workDir)[0] && 2 > strings.Count(workDir, "/") {
        die(fmt.Sprintf("No, it doesn’t works with root or any top-level directory. It will not process \"%s\" or any other \"/directory/\" under root.", workDir))
    }

    // Check, whether target "-dir" is a directory or a single file... And wether it’s exists at all...
    pathInfo, err := os.Stat(workDir)
    if err != nil {
        die(fmt.Sprintf("Error accessing the path \"%s\".\n", workDir)) // We could also output `err`, but it's confusing. Our simple message is enough.
    }

    expire, err := strconv.Atoi(*cliArgs["expire"].(*string))
    if nil != err {
        die(fmt.Sprintf("Invalid integer value in argument -expire=%s. Please use integer value to specify the number of days, or skip it to use default %d days.",
            *cliArgs["expire"].(*string), DefExpireDays))
    }

    // Calculate the expiration time
    expireDuration := time.Duration(expire) * 24 * time.Hour
    expireTime := time.Now().Add(-expireDuration)

    if pathInfo.IsDir() {
        // Walk through the directory and find files older than expireTime
        err = filepath.Walk(workDir, func(path string, pathInfo os.FileInfo, err error) error {
            if nil != err {
                return err
            }

            // Skip subdirectories if `-sub` is not specified
            isDir := pathInfo.IsDir()
            if isDir {
                // skip directory itself, we need its contents and sub-dirs.
                if workDir == path {
                    return nil
                }

                // skip subdirectories if we don't need them
                if !*cliArgs["sub"].(*bool) {
                    return filepath.SkipDir
                }
            }

            if !isDir {
                checkGarbageFile(path, pathInfo, expireTime) // Check if the file is older than expireTime
            }

            return nil // success
        })

        if nil != err {
            die(fmt.Sprintf("Error walking the path `%s`.\n%v\n", workDir, err))
        }

    // If -dir="[single file]"
    }else {
        checkGarbageFile(workDir, pathInfo, expireTime)
    }

    // Success
    fmt.Println("Done.")
}
