$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o "wordle-$env:GOOS-$env:GOARCH.exe"
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o "wordle-$env:GOOS-$env:GOARCH"
$env:GOOS="darwin"; $env:GOARCH="arm"; go build -o "wordle-$env:GOOS-$env:GOARCH"
$env:GOOS="freebsd"; $env:GOARCH="amd64"; go build -o "wordle-$env:GOOS-$env:GOARCH"
$env:GOOS="freebsd"; $env:GOARCH="arm"; go build -o "wordle-$env:GOOS-$env:GOARCH"
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o "wordle-$env:GOOS-$env:GOARCH"
$env:GOOS="linux"; $env:GOARCH="arm"; go build -o "wordle-$env:GOOS-$env:GOARCH"
$env:GOOS="android"; $env:GOARCH="arm"; go build -o "wordle-$env:GOOS-$env:GOARCH"