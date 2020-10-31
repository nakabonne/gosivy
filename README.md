# gosivy
Diagnose Go processes and visualize the results in real-time.

<demo>

Gosivy lets you visualize and diagnose your Go application no matter where it's running on.

## Usage
### Quickstart

Run the example application:
```
git clone https://github.com/nakabonne/gosivy.git
go run gosivy/examples/local/main.go
```

Then simply perform `gosivy` with no arguments (it automatically finds the process where the agent runs on):
```
gosivy
```

### Local mode
To diagnose a Go process running locally, launch agent as:

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
	time.Sleep(time.Hour)
}
```

```
gosivy -l
```

```go
gosivy 3400
```

Note that you need to start the `gosivy` process as the same user as the target application.

### Remote mode