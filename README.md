# XPB
A [pocketbase](https://pocketbase.io/) builder inspired by [xcaddy](https://github.com/caddyserver/xcaddy/) and [xpocketbase](https://github.com/kennethklee/xpb).

## Installation
### Dependencies
- [Go 1.23+](https://go.dev/doc/install)
### Download a Release
- Download the latest xpb release [here](https://github.com/PocketBuilds/xpb/releases)
### Use `go install`
1. Install [Go](https://go.dev/doc/install).
2. Follow [this guide](https://go.dev/doc/tutorial/compile-install) to add Go install directory to your system's shell path.
3. Run the following command:
```
go install github.com/PocketBuilds/xpb/cmd/xpb@latest
```
4. Optionally to get the version command to properly work you need to add an ldflag to the command with the specified version you want:
```
go install github.com/PocketBuilds/xpb/cmd@<version> -ldflags '-X github.com/PocketBuilds/xpb.version=<version>'
```
## Using the builder
```
$ xpb build [<pocketbase_version>]
    [--output <file>]
    [--with <module[@version][=replacement]>...]
    [--config <file>]
```
- `<pocketbase_version>` is the desired version of pocketbase to use. It defaults to the value of the env variable `XPB__PB_VERSION`, or to `latest` if that is not defined.
- `--output` is the output file. Defaults to `pocketbase`.
- `--with` is used to add an xpb plugin to the build. This flag can be used multiple times to add multiple plugins. It can also be used to change the version of xpb that is used for the build or replace the xpb and pocketbase modules to forked repositories if desired.
- `--config` should be a filepath to a `.toml` config file that describes your build. This can be used as an alternative to the cli arguments, or it can be used in conjunction with them. Config options described with cli arguments will take precedence and override config options in the config file if applicable.

## Creating a Plugin
You can use the following command to create a plugin from project in the current directory:
```
xpb plugin init <plugin_name>
```
- `<plugin_name>` is the unique name and identifier for your plugin. The command will also run `go mod init <plugin_name` assuming the go mod will be the same name as your plugin, but this can be manually changed in `go.mod` to match your git repo.

An XPB plugin is just a struct that adheres to the `xpb.Plugin` interface
