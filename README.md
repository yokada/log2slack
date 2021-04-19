# Local Development

```shell
go mod init log2slack
```

# Build

```shell
go build -o bin/log2slack cmd/log2slack.go
```

# Run

```shell
./bin/log2slack /path/to/file.log
```

# Test

```shell
go test ./... -v
```

# Refs

[ Upstart script for Go service ]( https://gist.github.com/sdrew/8e200bad0ce625f64c6d )

To help using inotify get through tail command in main function:
[ coreutils/tail.c at master · coreutils/coreutils ]( https://github.com/coreutils/coreutils/blob/master/src/tail.c )
