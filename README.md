# Garbage Collector

CLI tool to delete outdated files, expired for more than N days.

PLEASE BE CAREFUL! AUTHORS ARE NOT RESPONSIBLE IF YOU ACCIDENTALLY DELETE SOMETHING IMPORTANT!

# Install instructions for n00bs

1. Install GO. (Ask ChatGPT how to do this. In some cases you need to remove outdated version to install the fresh one.)
Usually it's something like
<pre>wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz</pre>
...then add GO to PATH, if needed, with
<pre>echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile</pre>
...make sure that you’re using fresh Go version with `go version`. If you see outdated version number, locate it with `which go` or `whereis go`, then remove duplicates with `sudo rm -rf /[path]/go` (BE CAREFUL!).

2. Run Garbage Collector with `go run garbage-collector.go` (follow CLI instructions).
3. Compile Garbage Collector with `go build garbage-collector.go`.
