# Description

A library that listens for commands sent either via HTTP or TCP. Commands are sent to an internal channel from where they are forwarded to registered listeners. Allows unlimited number of listeners that each receive the same command once.

# Usage

After importing the library in your program

```
test1.Listen(func(cmd test1.Cmd) { fmt.Println(1, cmd) })
test1.Listen(func(cmd test1.Cmd) { fmt.Println(2, cmd) })
# starts the http and tcp servers
# first parameter is the port for the http server
# second parameter is the port for the tcp server
test1.SetupServers(":9999", ":8888")
```