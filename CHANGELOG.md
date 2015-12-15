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
