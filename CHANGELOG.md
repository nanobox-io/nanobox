## v0.18.4 (March 9, 2016)

FEATURES:
  - Updated to use the new Mist (1.0.0)

## v0.18.3 (February 22, 2016)

IMPROVEMENTS:
  - Better error messaging when nanobox is unable to communicate with nanobox-server


## v0.18.2 (February 8, 2016)

BUG FIXES:
  - Nanobox create/destroy will now call the correct command when attempting to
  execute sudoed commands; this was caused after the move to dev subcommands.

## v0.18.1 (February 4, 2016)

FEATURES:
  - Nanobox will no longer create, publish, or fetch engines or services.

## v0.18.0 (February 4, 2016)

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

## v0.17.4 ()

  - Removed anything related to logtap since it isn't used in this way anymore
  (nanobox gets all its historical logs from a /logs route to the server)

## v0.17.3 (December 24, 2015)

BUG FIXES:
  - Fixed a regression caused by the previous version in which some clients were
  being prematurely closed causing panics when they were later attempted to close
  because it was presumed they were left open.

## v0.17.2 (December 22, 2015)

FEATURES:
  - Adds the ability to forward local proxy variables to the vm, for docker and
  nanobox server use (97fbfcc).

## v0.17.1 (December 21, 2015)

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

## v0.17.0 (December 21, 2015)

FEATURES:

  - Moved the update functionality into it's own script; the reason for this is
  to keep updating "safe". In the past there has been a break to the updating
  that has made it impossible to update to new versions of nanobox that contain
  fixes, meaning the only option is to manually download a new nanobox. Now however,
  the update script will work independently of nanobox, and nanobox will download
  and use the script when updating. If a break should occur the script can be
  run manually to update nanobox for fixes (30881d9).

## v0.16.17 (December 17, 2015)

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

## v0.16.16 (December 15, 2015)

FEATURES:

  - Added `dev_config` to the .nanofile parser to allow the setting of the guest
  vm environment (c4bf44a).
  - Added `dev-config` flag to `nanobox dev` allowing on-the-fly setting of the
  guest vm environment(c4bf44a).

## v0.16.15 (December 9, 2015)

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
