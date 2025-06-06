# Ember Mug Bluetooth Documentation

## Introduction

This repository contains a reverse-engineered documentation for the bluetooth API of [Ember Mugs](https://ember.com/).

The information provided here was retrieved using the Ember smartphone app on Android using an **Ember Mug 2** and decompiling the APK through [Java decompilers](http://www.javadecompilers.com/apk).

It may not be applicable to other Ember mugs, but feel free to extend the documentation and open a pull request :).

## Privacy

Data collected by the bluetooth gets sent to: [`https://collector.embertech.com`](https://collector.embertech.com)

Read [Data Collection & Privacy](./data-collection.md) for more information.

## Documentation

All commands have a service UUID of `fc543622236c4c948fa9944a3e5353fa`

* [Mug color](./docs/mug-color.md)
* [Target temperature](./docs/target-temp.md)
* [Current temperature](./docs/current-temp.md)
* [Battery percentage](./docs/battery.md)
* [Temperature unit](./docs/temperature-unit.md)
* [Liquid level](./docs/liquid-level.md)
* [Liquid state](./docs/liquid-state.md)
* [Mug name](./docs/mug-name.md)
* [Date & timezone](./docs/time-date-zone.md)
* [Push events](./docs/push-events.md)
* [Firmware & hardware versions](./docs/push-events.md)


## CLI Tool

A small command-line program is included to interact with an Ember Mug on Linux and other platforms. It now uses the Go Bluetooth library `tinygo.org/x/bluetooth`, so no external utilities are required.

### Build

```bash
$ go build ./cmd/embermug
```

You can store your mug's MAC address in `~/.config/embermug.json` to avoid
passing it to every command:

```json
{
  "mac": "AA:BB:CC:DD:EE:FF"
}
```

### Example

```bash
$ ./embermug status        # add --mac <MAC> if no config file
```

Set the target temperature to 55°C:

```bash
$ ./embermug set-target-temp --temp 55   # add --mac <MAC> if no config file
```

This will print the current temperature, target temperature and battery percentage.
