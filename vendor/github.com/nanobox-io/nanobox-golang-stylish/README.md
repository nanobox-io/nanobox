Nanobox Stylish
---------------

A stylish little library that styles output according to the nanobox style guide (http://nanodocs.gopagoda.io/engines/style-guide/)


### Installation

Frist `go get` the project...

`go get github.com/nanobox-io/nanobox-golang-stylish`


...then include the library in your project

`include github.com/nanobox-io/nanobox-golang-stylish`


### Usage

Available styles:
+ ProcessStart(process string)
+ ProcessEnd(process string)
+ Bullet(bullet string/[]string{"bullet", "bullet"...})
+ Warning(body string)
+ Fatal(header, body string)

For detailed information see the [complete documentation](http://godoc.org/github.com/nanobox-io/nanobox-golang-stylish)


### Examples

#### Processes
```go
stylish.ProcessStart("i am a process")
// process output
stylish.ProcessEnd()

// outputs
+ I am a process ------------------------------------------------------------ >
// process output
```

#### Bullets

##### single bullet
```go
stylish.Bullet("i am a bullet")

// outputs
+ i am a bullet
```

##### multiple bullets
```go
stylish.Bullet([]string{"we", "are", "many", "bullets"})

// outputs
+ we
+ are
+ many
+ bullets
```

#### Warnings
```go
stylish.Warning("i am a warning")

// outputs
----------------------------------  WARNING  ----------------------------------
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

Complete documentation is available on [godoc](http://godoc.org/github.com/nanobox-io/nanobox-golang-stylish).
