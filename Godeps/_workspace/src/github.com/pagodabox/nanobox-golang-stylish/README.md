Nanobox Stylish
---------------

A stylish little library that styles output according to the nanobox style guide (http://nanodocs.gopagoda.io/engines/style-guide/)


### Installation

Frist `go get` the project...

`go get github.com/pagodabox/nanobox-golang-stylish`


...then include the library in your project

`include github.com/pagodabox/nanobox-golang-stylish`


### Usage

Available styles:
+ Header(header string)
+ ProcessStart(process string)
+ ProcessEnd(process string)
+ SubTask(taks string)
+ SubTaskSuccess()
+ SubTaskFail()
+ Bullet(bullet string/[]string{"bullet", "bullet"...})
+ Warning(body string)
+ Fatal(header, body string)

For detailed information see the [complete documentation](http://godoc.org/github.com/pagodabox/nanobox-golang-stylish)


### Examples


#### Headers
```go
stylish.Header("i am a header")

// outputs
:::::::::::::::::::::::::: I AM A HEADER :::::::::::::::::::::::::
```


#### Processes
```go
stylish.ProcessStart("i am a process")
// process output
stylish.ProcessEnd("i am a process")

// outputs
I AM A PROCESS :::::::::::::::::::::::::::::::::::::::::::::::: =>
// process output
<= :::::::::::::::::::::::::::::::::::::::::::: END I AM A PROCESS
```


#### SubTasks

##### successful subtask
```go
stylish.SubTask("i am a successful sub task")
// subtask output
stylish.SubTaskSuccess()

// outputs
::::::::: I AM A SUCCESSFUL SUB TASK
// subtask output
<<<<<<<<< [√] SUCCESS
```

##### failed subtask
```go
stylish.SubTask("i am a failed sub task")
// subtask output
stylish.SubTaskFail()

// outputs
::::::::: I AM A FAILED SUB TASK
// subtask output
<<<<<<<<< [!] FAILED
```


#### Bullets

##### single bullet
```go
stylish.Bullet("i am a bullet")

// outputs
+> i am a bullet
```

##### multiple bullets
```go
stylish.Bullet([]string{"we", "are", "many", "bullets"})

// outputs
+> we
+> are
+> many
+> bullets
```

#### Warnings
```go
stylish.Warning("i am a warning")

// outputs
-----------------------------  WARNING  -----------------------------
i am a warning
```

#### Errors
```go
stylish.Error("i am an error", "things are probably going to explode now")

// outputs
! I AM AN ERROR !

things are probably going to explode now
```

### Documentation

Complete documentation is available on [godoc](http://godoc.org/github.com/pagodabox/nanobox-golang-stylish).
