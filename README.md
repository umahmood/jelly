# Jelly

Jelly is a back to basic no frills logger for Go applications. It is very simple 
to use, it provides three simple log levels. These are:

- jelly.Info - to log information *users* need to know about.
- jelly.Debug - to log information *developers* need to know about.
- jelly.Die - to log information about unrecoverable application states. 

The log files are stored in the users home directory under the *.jelly* directory.

The log format is as follows:

> DATE TIME UTC FILE_NAME:LINE_NO LOG_LEVEL - MESSAGE
 
For example:

> 2016/02/10 18:11:41 network.go:26 DEBUG - failed to create new node


# Installation

> go get github.com/umahmood/jelly

# Usage

A simple example:

    package main

    import (
        "fmt"
        "errors"
        
        "github.com/umahmood/jelly"
    )

    func main() {
        log, err := jelly.NewLog("app.log")
        if err != nil {
            fmt.Println(err)
        }

        fmt.Println("Log path", log.Path) // ~/.jelly/app.log
        fmt.Println("Log name", log.Name) // app.log

        log.Info("hello", "world", 42, true)
        
        err = errors.New("cruel world")
        log.Debug(err)
        
        // causes the application to exit with status code 1
        log.Die("he's dead jim")
    }

Sample output:

    > cat ~/.jelly/app.log
    2016/02/10 19:50:19 main:19 INFO - hello world 42
    2016/02/10 19:50:19 main.go:22 DEBUG - cruel world
    2016/02/10 19:50:19 main.go:25 DIE - he's dead jim

A single logger instance can be used across multiple goroutines:

    package main

    import (
        "fmt"
        "sync"

        "github.com/umahmood/jelly"
    )

    func main() {
        log, err := jelly.NewLog("app.log")
        if err != nil {
            fmt.Println(err)
        }

        var wg sync.WaitGroup

        work := func(items []string, b bool, log *jelly.Logger) {
            defer wg.Done()
            for _, i := range items {
                if b {
                    log.Debug(i)
                } else {
                    log.Info(i)
                }
            }
        }

        items := []string{"A", "B", "C", "D", "E"}
        wg.Add(1)
        go work(items, true, log)

        items = []string{"F", "G", "H", "I", "J"}
        wg.Add(1)
        go work(items, false, log)

        wg.Wait()
    }

Sample output:

    > cat ~/.jelly/app.log
    2016/02/10 19:58:01 main.go:24 INFO - F
    2016/02/10 19:58:01 main.go:24 INFO - G
    2016/02/10 19:58:01 main.go:22 DEBUG - A
    2016/02/10 19:58:01 main.go:24 INFO - H
    2016/02/10 19:58:01 main.go:22 DEBUG - B
    2016/02/10 19:58:01 main.go:24 INFO - I
    2016/02/10 19:58:01 main.go:22 DEBUG - C
    2016/02/10 19:58:01 main.go:24 INFO - J
    2016/02/10 19:58:01 main.go:22 DEBUG - D
    2016/02/10 19:58:01 main.go:22 DEBUG - E

# Documentation

> http://godoc.org/github.com/umahmood/jelly

# License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
