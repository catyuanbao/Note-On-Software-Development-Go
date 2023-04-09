## About ReadLine()

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    line, isPrefix, err := reader.ReadLine()
    if err != nil {
        fmt.Println("Error reading input:", err)
        return
    }
    if isPrefix {
        fmt.Println("Error: input line too long")
        return
    }
    fmt.Println("Input:", string(line))
}

```

```c
#include <stdio.h>

int main() {
        printf("BUFSIZ: %d\n", BUFSIZ);
        return 0;
}
```

output is:

```
BUFSIZ: 8192
```

if read line more than 8192, got an error

```bash
>>> go run main.go < bigfile                                                                                                        ‹git:main ✘› 16:35.03 Sun Apr 09 2023 >>> 
Error: input line too long
```

