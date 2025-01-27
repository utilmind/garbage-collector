# Garbage Collector in Go language

CLI tool to remove from specified directory outdated/expired files, created more than N days ago.

PLEASE BE CAREFUL! AUTHORS ARE NOT RESPONSIBLE IF YOU ACCIDENTALLY DELETE SOMETHING IMPORTANT!

<blockquote>UPD. This tool was created because author was not aware that Linux already have standard tool `find`. 🤡
  
E.g. `$ find <path> -type f -mtime <days> -delete`.
However you still can use this code if you prefer to have cross-platform code for this purpose, or would like to use it as the base for your customized solution.</blockquote>

# TODO
* Support multiple paths in single run, as `-dir=[...]` argument. Either comma `,` or pipe `|` separated.
* Support search depth (similar to the depth of Linux `find` tool).

# Install instructions for noobs

1. Install GO. (In some cases you need to remove outdated version to install the fresh one.)
Usually it's something like `sudo apt install golang-go`, or if you need the newest version, then check the fresh version number on https://go.dev/doc/install and use something like the following:
<pre>wget https://go.dev/dl/go1.23.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.5.linux-amd64.tar.gz</pre>
...then add GO to PATH, if needed, with
<pre>echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile</pre>
...make sure that you’re using fresh Go version with `go version`. If you see outdated version number, locate it with `which go` or `whereis go`, then remove duplicates with `sudo rm -rf /[path]/go` (BE CAREFUL!).

2. Run Garbage Collector with `go run garbage-collector.go` (follow CLI instructions).
3. Compile Garbage Collector with `go build garbage-collector.go`.
