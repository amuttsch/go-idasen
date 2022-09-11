# go-idasen

![Go](https://github.com/amuttsch/go-idasen/workflows/Go/badge.svg)

This is a Go adaptation of [newAM/idasen](https://github.com/newAM/idasen) which has a similar
command line interface and configuration file format.

The IDÅSEN is an electric sitting standing desk with a Linak controller sold by
ikea.

The position of the desk can controlled by a physical switch on the desk or
via bluetooth using an phone app.

This is a command line interface and API written in Go to control the Idasen via
bluetooth from a desktop computer.

## Set Up

### Prerequisites

The desk should be connected and paired to the computer.

### Install

```shell
$ go get -u github.com/amuttsch/go-idasen
```

## Configuration
Configuration that is not expected to change frequency can be provided via a
YAML configuration file located at ``./idasen.yaml``.

You can use this command to initialize a new configuartion file:

```shell
$ go-idasen init
```


```yaml
mac_address: AA:AA:AA:AA:AA:AA
positions:
    sit: 0.75
    stand: 1.1
```

Configuration options:

* ``mac_address`` - The MAC address of the desk. This is required.
* ``positions`` - A dictionary of positions with values of desk height from the
  floor in meters, ``sit`` and ``stand`` are provided as examples.

The program will try to find the device address,
but if it fails, it has to be done manually.

The device MAC addresses can be found using ``blueoothctl`` and bluetooth
adapter names can be found with ``hcitool dev`` on linux.

## Usage

### Command Line

To print the current desk height:

    $ go-idasen height

To monitor for changes to height:

    $ go-idasen monitor

To save the current height as the sitting position:

    $ go-idasen save sit

To delete the saved sitting position:

    $ go-idasen delete sit

To move the desk to a specific height:

    $ go-idasen move 0.8

Assuming the config file is populated to move the desk to sitting position:

    $ go-idasen sit

### API

Check `cmd/` for examples on how to use the API.

Import `go-idasen` in your application:

    import "github.com/amuttsch/go-idasen/idasen"

Except for desk discovery you should always defer `idasen.Close()` to clean up and disconnect the bluetooth connection.

#### Discover a desk

To discover a desk:

    desk, err := idasen.DiscoverDesk()

It returns a `desk` struct containing the name and mac address of the desk or an error if no desk was found.

#### Current height

To get the current height:

```go
idasen, err := idasen.New(config)
if err != nil {
    fmt.Println(err)
    return
}
defer idasen.Close()

h, err := idasen.HeightInMeters()
if err != nil {
    fmt.Println(err)
    return
}
```

#### Move the desk

First initialize the desk:

```go
idasen, err := idasen.New(config)
if err != nil {
    fmt.Println(err)
    return
}
defer idasen.Close()
```

Then you can either move the desk in one direction, stop movement or move it to a target position:

```go
idasen.MoveUp()
idasen.MoveDown()
idasen.MoveStop()
idasen.MoveToTarget(0.9) // in meters
```
