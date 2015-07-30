## Hatchet

Hatchet is a simple Logger interface that is intentionally generic which allows it to be used by numerous existing loggers, or custom loggers.
It also provides a DevNullLogger which can be used in projects that require loggers and none is provided/desired


### Usage

Hatchet is best suited for projects that have a core package that is importing numerous external packages which need loggers:

    // package main
    package main

      import(
        "github.com/jcelliott/lumber"
        "github.com/nanobox-core/mist"
      )

      //
      func main() {

        // lumber is merely an example of an existing logging package that satsifies
        // hatchets Logger interface.
        logger := lumber.NewConsoleLogger(lumber.DEBUG)

        // mist requires a logger
        mist := mist.New(logger)
      }

    // package mist
    package mist
      import "github.com/nanobox-core/hatchet"

      // mist requires a logger
      type Mist struct {
        log hatchet.Logger
      }

      // New creates a new mist, and sets its logger
      func New(logger hatchet.Logger) *Mist {

        // if no logger is provided, hatchets DevNullLogger is used.
        if logger == nil {
          logger = hatchet.DevNullLogger{}
        }

        mist := &Mist{
          log: logger
        }

        mist.Log.Info("Created new mist...\n")

        return mist
      }


### Documentation

Complete documentation is available on [godoc](http://godoc.org/github.com/nanobox-core/hatchet).


### Contributing
