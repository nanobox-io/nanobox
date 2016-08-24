package code

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  
  "github.com/nanobox-io/nanobox/generators/hooks/build"
  "github.com/nanobox-io/nanobox/util/hookit"
)

// RunUserHook runs the user hook inside of the specified container
func RunUserHook(container string) (string, error) {
  // generate the user payload
  userPayload, err := build.UserPayload()
  if err != nil {
    lumber.Error("code:RunUserHook:build.UserPayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for user hook: %s", err.Error())
  }
  
  // run the user hook
  res, err := hookit.Exec(container, "user", userPayload, "debug")
  if err != nil {
    lumber.Error("code:RunUserHook:hookit.Exec(%s, %s, %s, %s): %s", container, "user", userPayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute user hook: %s", err.Error())
  }
  
  return res, nil
}

// RunConfigureHook runs the configure hook inside of the specified container
func RunConfigureHook(container string) (string, error) {
  // generate the configure payload
  configurePayload, err := build.ConfigurePayload()
  if err != nil {
    lumber.Error("code:RunConfigureHook:build.ConfigurePayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for configure hook: %s", err.Error())
  }
  
  // run the configure hook
  res, err := hookit.Exec(container, "configure", configurePayload, "debug")
  if err != nil {
    lumber.Error("code:RunConfigureHook:hookit.Exec(%s, %s, %s, %s): %s", container, "configure", configurePayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute configure hook: %s", err.Error())
  }
  
  return res, nil
}

// RunFetchHook runs the fetch hook inside of the specified container
func RunFetchHook(container string) (string, error) {
  // generate the fetch payload
  fetchPayload, err := build.FetchPayload()
  if err != nil {
    lumber.Error("code:RunFetchHook:build.FetchPayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for fetch hook: %s", err.Error())
  }
  
  // run the fetch hook
  res, err := hookit.Exec(container, "fetch", fetchPayload, "debug")
  if err != nil {
    lumber.Error("code:RunFetchHook:hookit.Exec(%s, %s, %s, %s): %s", container, "fetch", fetchPayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute fetch hook: %s", err.Error())
  }
  
  return res, nil
}

// RunSetupHook runs the setup hook inside of the specified container
func RunSetupHook(container string) (string, error) {
  // generate the setup payload
  setupPayload, err := build.SetupPayload()
  if err != nil {
    lumber.Error("code:RunSetupHook:build.SetupPayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for setup hook: %s", err.Error())
  }
  
  // run the setup hook
  res, err := hookit.Exec(container, "setup", setupPayload, "debug")
  if err != nil {
    lumber.Error("code:RunSetupHook:hookit.Exec(%s, %s, %s, %s): %s", container, "setup", setupPayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute setup hook: %s", err.Error())
  }
  
  return res, nil
}

// RunBoxfileHook runs the boxfile hook inside of the specified container
func RunBoxfileHook(container string) (string, error) {
  // generate the boxfile payload
  boxfilePayload, err := build.BoxfilePayload()
  if err != nil {
    lumber.Error("code:RunBoxfileHook:build.BoxfilePayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for boxfile hook: %s", err.Error())
  }
  
  // run the boxfile hook
  res, err := hookit.Exec(container, "boxfile", boxfilePayload, "debug")
  if err != nil {
    lumber.Error("code:RunBoxfileHook:hookit.Exec(%s, %s, %s, %s): %s", container, "boxfile", boxfilePayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute boxfile hook: %s", err.Error())
  }
  
  return res, nil
}

// RunPrepareHook runs the prepare hook inside of the specified container
func RunPrepareHook(container string) (string, error) {
  // generate the prepare payload
  preparePayload, err := build.PreparePayload()
  if err != nil {
    lumber.Error("code:RunPrepareHook:build.PreparePayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for prepare hook: %s", err.Error())
  }
  
  // run the prepare hook
  res, err := hookit.Exec(container, "prepare", preparePayload, "debug")
  if err != nil {
    lumber.Error("code:RunPrepareHook:hookit.Exec(%s, %s, %s, %s): %s", container, "prepare", preparePayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute prepare hook: %s", err.Error())
  }
  
  return res, nil
}

// RunCompileHook runs the compile hook inside of the specified container
func RunCompileHook(container string) (string, error) {
  // generate the compile payload
  compilePayload, err := build.CompilePayload()
  if err != nil {
    lumber.Error("code:RunCompileHook:build.CompilePayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for compile hook: %s", err.Error())
  }
  
  // run the compile hook
  res, err := hookit.Exec(container, "compile", compilePayload, "debug")
  if err != nil {
    lumber.Error("code:RunCompileHook:hookit.Exec(%s, %s, %s, %s): %s", container, "compile", compilePayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute compile hook: %s", err.Error())
  }
  
  return res, nil
}

// RunPackAppHook runs the pack-app hook inside of the specified container
func RunPackAppHook(container string) (string, error) {
  // generate the pack-app payload
  packAppPayload, err := build.PackAppPayload()
  if err != nil {
    lumber.Error("code:RunPackAppHook:build.PackAppPayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for pack-app hook: %s", err.Error())
  }
  
  // run the pack-app hook
  res, err := hookit.Exec(container, "pack-app", packAppPayload, "debug")
  if err != nil {
    lumber.Error("code:RunPackAppHook:hookit.Exec(%s, %s, %s, %s): %s", container, "pack-app", packAppPayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute pack-app hook: %s", err.Error())
  }
  
  return res, nil
}

// RunPackBuildHook runs the pack-build hook inside of the specified container
func RunPackBuildHook(container string) (string, error) {
  // generate the pack-build payload
  packBuildPayload, err := build.PackBuildPayload()
  if err != nil {
    lumber.Error("code:RunPackBuildHook:build.PackBuildPayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for pack-build hook: %s", err.Error())
  }
  
  // run the pack-build hook
  res, err := hookit.Exec(container, "pack-build", packBuildPayload, "debug")
  if err != nil {
    lumber.Error("code:RunPackBuildHook:hookit.Exec(%s, %s, %s, %s): %s", container, "pack-build", packBuildPayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute pack-build hook: %s", err.Error())
  }
  
  return res, nil
}

// RunCleanHook runs the clean hook inside of the specified container
func RunCleanHook(container string) (string, error) {
  // generate the clean payload
  cleanPayload, err := build.CleanPayload()
  if err != nil {
    lumber.Error("code:RunCleanHook:build.CleanPayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for clean hook: %s", err.Error())
  }
  
  // run the clean hook
  res, err := hookit.Exec(container, "clean", cleanPayload, "debug")
  if err != nil {
    lumber.Error("code:RunCleanHook:hookit.Exec(%s, %s, %s, %s): %s", container, "clean", cleanPayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute clean hook: %s", err.Error())
  }
  
  return res, nil
}

// RunPackDeployHook runs the pack-deploy hook inside of the specified container
func RunPackDeployHook(container string) (string, error) {
  // generate the pack-deploy payload
  packDeployPayload, err := build.PackDeployPayload()
  if err != nil {
    lumber.Error("code:RunPackDeployHook:deploy.PackBuildPayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for pack-deploy hook: %s", err.Error())
  }
  
  // run the pack-deploy hook
  res, err := hookit.Exec(container, "pack-deploy", packDeployPayload, "debug")
  if err != nil {
    lumber.Error("code:RunPackDeployHook:hookit.Exec(%s, %s, %s, %s): %s", container, "pack-deploy", packDeployPayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute pack-deploy hook: %s", err.Error())
  }
  
  return res, nil
}
