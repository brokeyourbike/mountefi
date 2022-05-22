# mountefi

[![Maintainability](https://api.codeclimate.com/v1/badges/94e83790ba593e90029a/maintainability)](https://codeclimate.com/github/brokeyourbike/mountefi/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/94e83790ba593e90029a/test_coverage)](https://codeclimate.com/github/brokeyourbike/mountefi/test_coverage)

EFI Mounting utility with no dependencies written in Golang.

## How it's working?

Instead of parsing unparsable plist returned from the the `diskutil list -plist` command like the original [MountEFI](https://github.com/corpnewt/MountEFI) is doing, we fetching only `AllDisks` property, and then calling `diskutil info -plist <disk>` for each of the disks concurrently.

Of course this approach will be a bit slower, but `diskutil info -plist <disk>` returning much more predictable plist format.