# gosivy

<demo>

Gosivy lets you diagnose the Go process no matter where it's running on, and visualizes the statistics in real-time. It helps you understand how your application consumes the resources.

## Installation
Binary releases are available through [here](https://github.com/nakabonne/gosivy/releases).

#### MacOS

```
brew install nakabonne/gosivy/gosivy
```

#### RHEL/CentOS

```
rpm -ivh https://github.com/nakabonne/gosivy/releases/download/v0.1.0/gosivy_0.1.0_linux_amd64.rpm
```

#### Debian/Ubuntu

```
wget https://github.com/nakabonne/gosivy/releases/download/v0.1.0/gosivy_0.1.0_linux_amd64.deb
apt install ./gosivy_0.1.0_linux_amd64.deb
```

#### Go
Required: >= 1.10. Note that you may have a problem because it downloads an untagged binary.
```
go get github.com/nakabonne/gosivy
```

#### Docker

```
docker run --rm -it nakabonne/gosivy gosivy
```

## Quickstart

Run the example application:
```
git clone https://github.com/nakabonne/gosivy.git
go run gosivy/examples/local/main.go
```

Then simply perform `gosivy` with no arguments (it automatically finds the process where the agent runs on):
```
gosivy
```

## Usage
First up, you start the agent in the process where you want to collect statistics. Then execute `gosivy` to scrape from the agent periodically. You can diagnose processes running not only locally (local mode), but also on another host (remote mode).

### Local Mode
Import `github.com/nakabonne/gosivy/agent` into the target application to launch the agent. See the [local example](https://github.com/nakabonne/gosivy/blob/main/examples/local/main.go).

```go
package main

import (
	"log"
	"time"

	"github.com/nakabonne/gosivy/agent"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	defer agent.Close()
}
```

With `-l` flag can list the processes where the agent runs on:
```
gosivy -l
```

Give the PID of the process to be diagnosed:
```
gosivy 5000
```

Be sure to start the `gosivy` process as the same user as the target application.

### Remote Mode
```go
package main

import (
	"log"
	"time"

	"github.com/nakabonne/gosivy/agent"
)

func main() {
    // TODO: Start agent
}
```

Give the target's agent address:

```
gosivy host.xz:8080
```

### Settings
Command-line options are:

```
Usage:
  gosivy [flags] <pid|host:port>

Flags:
      --debug                      Run in debug mode.
  -l, --list-processes             Show processes where gosivy agent runs on.
      --scrape-interval duration   Interval to scrape from the agent. It must be >= 1s (default 1s)
  -v, --version                    Print the current version.
```

Gosivy requires the config directory for pid management. By default it will be created undernearth `$HOME/.config` (`APPDATA` on windows). For those who want to assign another directory, `GOSIVY_CONFIG_DIR` environment variable is available.

## Supported Architectures
The architectures the agent can perform on are:

- FreeBSD i386/amd64/arm
- Linux i386/amd64/arm(raspberry pi)
- Windows/amd64
- Darwin i386/amd64
- OpenBSD amd64
- Solaris amd64

Those partially supported:
- CPU on DragonFly BSD
- host on Linux RISC-V

## Built With
- [shirou/gopsutil](https://github.com/shirou/gopsutil)
- [mum4k/termdash](https://github.com/mum4k/termdash)
  - [nsf/termbox-go](https://github.com/nsf/termbox-go)

## Inspired By
- [google/gops](https://github.com/google/gops)
- [arl/statsviz](https://github.com/arl/statsviz)


## Author
Ryo Nakao <<ryo@nakao.dev>>
