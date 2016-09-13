# go-util
A miscellaneous collection of utility functions for golang

## Installation
`go get github.com/sdomino/go-util`

## Usage
``` golang

  // It is recommended when importing, to declare the import as [package]util to avoid possible collisions
  import (
    cryptoutil "github.com/sdomino/go-util/crypto"
    fileutil "github.com/sdomino/go-util/file"
    printutil "github.com/sdomino/go-util/print"
    storageutil "github.com/sdomino/go-util/storage"
  )

  // example using printutil
  printutil.Color("Green [green]eggs[reset] and [red]ham[reset]")
```
