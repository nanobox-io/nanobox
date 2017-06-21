## 2.1.2 (June 1, 2017)

FEATURES:
  - Add team support for nanobox production commands [#453](https://github.com/nanobox-io/nanobox/pull/453)
  - Now much smarter about how it creates networks on native [#447](https://github.com/nanobox-io/nanobox/pull/447)

BUG FIXES:
  - Fix an idempotency problem on service creation [#454](https://github.com/nanobox-io/nanobox/pull/454)
  - Fix a permission denied problem on osx [#452](https://github.com/nanobox-io/nanobox/pull/452)
  - Add mount checking [#443](https://github.com/nanobox-io/nanobox/pull/443)
  - Make the ping system better at knowing when the server is up [#439](https://github.com/nanobox-io/nanobox/pull/439)

## 2.1.1 (May 15, 2017)

FEATURES:
  - Submit failures and logs to the production nanobox server to help coordinate with tickets [#436](https://github.com/nanobox-io/nanobox/pull/436)
  - Allow a file with evar values [#435](https://github.com/nanobox-io/nanobox/pull/435)
  - Adjust the way we do nfs config on windows [#430](https://github.com/nanobox-io/nanobox/pull/430)
  - Make tap a part of the server start on OSX [#428](https://github.com/nanobox-io/nanobox/pull/428)

BUG FIXES:
  - Fix a duplicate etc/exports issue on linux [#433](https://github.com/nanobox-io/nanobox/pull/433)
  - Remove a server shutdown issue when there is no server [#431](https://github.com/nanobox-io/nanobox/pull/431)

## 2.1.0 (April 7, 2017)

FEATURES:
  - Make a major change to the way nanobox runs. When running nanobox now creates a nanobox server
    This change fixes the quality of life of most users because it should only ask for passwords once
  - VPN now runs under the server
  - All administrative commands (dns add, sharing etc.) now run through the server


## 2.0.4 (March 6, 2017)

BUG FIXES:
  - Add idempotency to linux systemd service start [#399](https://github.com/nanobox-io/nanobox/pull/399)
  - Fix an issue where the dev container disappeard unexpectedly [#402](https://github.com/nanobox-io/nanobox/pull/402)
  - Fix an issue that caused us to pull s3 for updates every time [#396](https://github.com/nanobox-io/nanobox/pull/396)
  - Fix an issue with /etc/exports on osx [#392](https://github.com/nanobox-io/nanobox/pull/392)

FEATURES:
  - Make linux startup system more flexable [#397](https://github.com/nanobox-io/nanobox/pull/397)
  - Optimize display of summarized test [#395](https://github.com/nanobox-io/nanobox/pull/395)
  - Add a check in for boxfile existance [#388](https://github.com/nanobox-io/nanobox/pull/388)

## 2.0.3 (February 23, 2017)

BUG FIXES: 
  - If during the setup the VM errors we now clean up [#373](https://github.com/nanobox-io/nanobox/pull/373)
  - Detect the correct location for systemd on linux [#374](https://github.com/nanobox-io/nanobox/pull/374)
  - Fix an issue where docker being down would remove components [#376](https://github.com/nanobox-io/nanobox/pull/376)
  - Remove Duplicate paths in /etc/exports on osx [#379](https://github.com/nanobox-io/nanobox/pull/379)
  - Stop releasing ips during nanobox stop [#380](https://github.com/nanobox-io/nanobox/pull/380)

FEATURES:
  - Add messaging to help make it clear when networks fail [#378](https://github.com/nanobox-io/nanobox/pull/378)
  - Allow users to set the user they want to console in as [#385](https://github.com/nanobox-io/nanobox/pull/385)
  - Confirm docker connections during the init process [#386](https://github.com/nanobox-io/nanobox/pull/386)
  - Confirm VirtualMachine can talk to the host [#387](https://github.com/nanobox-io/nanobox/pull/387)


## 2.0.2 (February 14, 2017)

BUG FIXES: 
  - Allow vt100 terminal codes to work properly [#370](https://github.com/nanobox-io/nanobox/pull/370)
  - Fix Truncation of the summary to also work with headers [#369](https://github.com/nanobox-io/nanobox/pull/369)
  - Fix a bug where leftover [#365](https://github.com/nanobox-io/nanobox/issues/365)
  - Hook failures no longer show duplicates [#360s](https://github.com/nanobox-io/nanobox/issues/360s)

FEATURES:
  - In consoles use the exit code instead of our own[#361](https://github.com/nanobox-io/nanobox/issues/361)

## 2.0.1 (February 9, 2017)

BUG FIXES:
  - Remove duplicate display of error messages [#360](https://github.com/nanobox-io/nanobox/issues/360)
  - Allow hook timeouts to work properly

FEATURES:
  - Make our busgnag reporting much better [#353](https://github.com/nanobox-io/nanobox/issues/353)
  - Update the config command to be cleaner [#349](https://github.com/nanobox-io/nanobox/issues/349)
  - Adjust the error handling so stack tracers are cleaner
  - Allow for odin messages to be encapsulated inside the correct context when erroring
  - Add an error message when a evar is not added [#338](https://github.com/nanobox-io/nanobox/issues/338)

## 2.0.0 (January 31, 2017)

FEATURES:
  - Make nanobox work using docker-machine instead of vagrant
  - Rework networking to use a vpn
  - Add Native support to use native docker on the system

## 0.18.4 (March 9, 2016)

FEATURES:
  - Updated to use the new Mist (1.0.0)

## 0.18.3 (February 22, 2016)

IMPROVEMENTS:
  - Better error messaging when nanobox is unable to communicate with nanobox-server


## 0.18.2 (February 8, 2016)

BUG FIXES:
  - Nanobox create/destroy will now call the correct command when attempting to
  execute sudoed commands; this was caused after the move to dev subcommands.

## 0.18.1 (February 4, 2016)

FEATURES:
  - Nanobox will no longer create, publish, or fetch engines or services.

## 0.18.0 (February 4, 2016)

BREAKING CHANGES:
  - All relevant "dev" commands are now sub-commands of "dev" in preparation
  for production commands.

FEATURES:
  - Added the ability to create, publish, and fetch "services" similar to the way
  "engines" are now.
  - Removed "overlays" when publishing or fetching engines; engines now take a
  list of files required for the build, specified in the Enginefile (this is done
  as part of the most away from "zero config")

IMPROVEMENTS:
  - More tests added.
  - Removed all "default" and "common" stuff for now in favor of a cleaner
  implementation later

BUG FIXES:
  - Fixed an issue where the raw terminal was being closed early causing double
  output and improper capturing of signals connected via "terminal commands".
  - Fixed an issue causing deploys to happen on every dev, rather than only after
  a recent reload.

## 0.17.4 ()

  - Removed anything related to logtap since it isn't used in this way anymore
  (nanobox gets all its historical logs from a /logs route to the server)

## 0.17.3 (December 24, 2015)

BUG FIXES:
  - Fixed a regression caused by the previous version in which some clients were
  being prematurely closed causing panics when they were later attempted to close
  because it was presumed they were left open.

## 0.17.2 (December 22, 2015)

FEATURES:
  - Adds the ability to forward local proxy variables to the vm, for docker and
  nanobox server use (97fbfcc).

## 0.17.1 (December 21, 2015)

IMPROVEMENTS:

  - Moved a significant amount of logic related to creating/managing the tty
  terminal out of the server package and into the terminal package (61399a2).
  - Turned the execInternal command into three commands that would relate
  specifically to what they were doing (Console, Develop, Exec); this was done
  to help separate out functionality and provide clarity (61399a2).
  - Added a function (Connect) to mist that emulates the Stream functionality, but returns
  the client to be closed externally to the function (61399a2).
  - 'nanobox dev' will now always connect to mist (not just on deploy), passing
  the connection into its respective sever call, which closes the connection after
  it receives data from the server; this is done to allow hooks to contain output
  that can be displayed to the client via mist (61399a2).
  - More idiomatic way of checking flag in dev command (59b8735).

## 0.17.0 (December 21, 2015)

FEATURES:

  - Moved the update functionality into it's own script; the reason for this is
  to keep updating "safe". In the past there has been a break to the updating
  that has made it impossible to update to new versions of nanobox that contain
  fixes, meaning the only option is to manually download a new nanobox. Now however,
  the update script will work independently of nanobox, and nanobox will download
  and use the script when updating. If a break should occur the script can be
  run manually to update nanobox for fixes (30881d9).

## 0.16.17 (December 17, 2015)

IMPROVEMENTS:

  - Moved some file close defers to be more idiomatic (c615cab).
  - Added newlines to various log output where it had been missed (c615cab).
  - Use config rather than runtime to determine environment in update.go (c615cab).
  - Minor cleanup of merge #119 (c615cab).
  - Only show vagrant log error once (rather than per error) to reduce log spam (e7060e0).
  - Removed blank newline character from vagrant output (51d286f).
  - Show 100% when nanobox/boot2docker is done downloading (51d286f).

BUG FIXES:

  - Merged #119 - Check to ensure the newly downloaded CLI matches the remote md5
  (this also fixed issue #116).
  - Removed empty vagrant output lines (f076651).
  - Prompting for admin before delete runs to avoid password prompt showing up in
  the middle of vagrant output (f076651).

## 0.16.16 (December 15, 2015)

FEATURES:

  - Added `dev_config` to the .nanofile parser to allow the setting of the guest
  vm environment (c4bf44a).
  - Added `dev-config` flag to `nanobox dev` allowing on-the-fly setting of the
  guest vm environment(c4bf44a).

## 0.16.15 (December 9, 2015)

FEATURES:

  - **New Command: `nanobox reload`**: This will reload the nanobox by suspending
  the VM and resuming it. This is an effective way to attempt to recover a VM
  before destroying it completely (84e7d68).

IMPROVEMENTS:

  - Cleaned up some crufty code to improve readability (9b4e0d3).

BUG FIXES:

  - Merged #113 - Fixes the first read of when watching files.


## Previous (December 9, 2015)

This change log began with version 0.16.15. Any prior changes can be seen by viewing
the commit history.
