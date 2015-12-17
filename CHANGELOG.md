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

This changelog began with version 0.16.15. Any prior changes can be seen by viewing
the commit history.
