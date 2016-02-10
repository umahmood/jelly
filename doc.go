/*
Package jelly a no frills logger for Go applications.

Example usage:

    package main

    import (
        "fmt"

        "github.com/umahmood/jelly"
    )

    func main() {
        log, err := jelly.NewLog("app.log")
        if err != nil {
            // ...
        }

        //...

        log.Info("hello", "world", 42, true)

        // ...

        log.Debug("cruel", "world", 66.6, false)

        // ...

        log.Die("houston, we have a problem", err)
    }
*/
package jelly
