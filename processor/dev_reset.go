package processor

import (
  "github.com/nanobox-io/nanobox/util/counter"
  "github.com/nanobox-io/nanobox/util/data"
)

type devReset struct {
  config ProcessConfig
}

func init() {
  Register("dev_reset", devResetFunc)
}

func devResetFunc(config ProcessConfig) (Processor, error) {
  return devReset{config: config}, nil
}

func (self devReset) Results() ProcessConfig {
  return self.config
}

func (self devReset) Process() error {

  if err := self.resetCounters(); err != nil {
    return err
  }

  return nil
}

// resetCounters resets all the counters associated with all apps
func (self devReset) resetCounters() error {

  // reset the provider counter
  counter.Reset("provider")

  apps, err := data.Keys("apps")

  if err != nil {
    return err
  }

  for _, app := range apps {
    // reset the general app usage counter
    counter.Reset(app)

    // reset the app dev usage counter
    counter.Reset(app + "_dev")
  }

  return nil
}
