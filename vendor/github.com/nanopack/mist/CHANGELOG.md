## v1.1.2 (August 18, 2016)

IMPROVEMENTS:
  - Layer in WSS

## v1.1.1 (August 5, 2016)

BUG FIXES:
  - Resolve stray goroutine memory leak
  - Ensure error is returned if auth is unsuccessful

IMPROVEMENTS:
  - All connections are not subscribers
  - Allow server as a config file option
  - Add versioning
  - Enhance logging

## v1.1.0 (April 22, 2016)

BUG FIXES:
  - Authentication is now handled on a per connection basis, with each connection
  required to send a token if the server is started with an authenticator.
  - Subscriptions have been reworked to meet the original requirements of mist.
  - Fixed an issue with how mist was parsing config files.

IMPROVEMENTS:
  - Clients now have a way to authenticate with an authenticated server.
  - Response messages have been added when running commands against a server to
  verify that the command worked/failed.
  - Many tests have been added/updated.

## v1.0.2 (March 17, 2016)

BUG FIXES:
  - Updated how Mist parses configs.

IMPROVEMENTS:
  - Added additional options for specifying logging type (stdout, file) and an option
  for specifying the log path.

## v1.0.1 (March 10, 2016)

BUG FIXES:
  - Break if mist fails to write to a web socket client; it's presumably dead, and
  we don't want mist looping forever trying to read/write to something git will never
  be able to.

## Previous (March 9, 2016)

This change log began with version 1.0.0. Any prior changes can be seen by viewing
the commit history.
