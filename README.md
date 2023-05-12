# mountefi

[![Go Reference](https://pkg.go.dev/badge/github.com/brokeyourbike/mountefi.svg)](https://pkg.go.dev/github.com/brokeyourbike/mountefi)
[![Go Report Card](https://goreportcard.com/badge/github.com/brokeyourbike/mountefi)](https://goreportcard.com/report/github.com/brokeyourbike/mountefi)
[![Maintainability](https://api.codeclimate.com/v1/badges/94e83790ba593e90029a/maintainability)](https://codeclimate.com/github/brokeyourbike/mountefi/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/94e83790ba593e90029a/test_coverage)](https://codeclimate.com/github/brokeyourbike/mountefi/test_coverage)

EFI Mounting utility with no dependencies written in Golang.

## How it's working?

Instead of parsing unparsable plist returned from the `diskutil list -plist` command like the original [MountEFI](https://github.com/corpnewt/MountEFI) is doing, we fetching `AllDisks` property (contains a list of disks), and then calling `diskutil info -plist <disk>` for each of the disks concurrently.

Of course this approach can be a bit slower, but `diskutil info -plist <disk>` returning much more predictable plist format.

## Authors
- [Ivan Stasiuk](https://github.com/brokeyourbike) | [Twitter](https://twitter.com/brokeyourbike) | [LinkedIn](https://www.linkedin.com/in/brokeyourbike) | [stasi.uk](https://stasi.uk)

## Thanks
https://github.com/corpnewt/MountEFI

## License
[BSD-3-Clause License](https://github.com/brokeyourbike/mountefi/blob/main/LICENSE)
