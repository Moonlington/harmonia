# Harmonia

[![Go Reference](https://pkg.go.dev/badge/github.com/Moonlington/harmonia.svg)](https://pkg.go.dev/github.com/Moonlington/harmonia) [![Go Report Card](https://goreportcard.com/badge/github.com/Moonlington/harmonia)](https://goreportcard.com/report/github.com/Moonlington/harmonia)

**Harmonia** is a [Go](https://golang.org/) module that is built on top of [DiscordGo](https://github.com/bwmarrin/discordgo). It aims to make the developing process of a Go based Discord bot easier by giving developers the tools needed to effectively and quickly design bots.

Anything you can do in [DiscordGo](https://github.com/bwmarrin/discordgo) is doable in **Harmonia**.

***This module is still in heavy development!***

## Getting Started with Harmonia

### Installing

Assuming you already have a working Go environment, run the following in your project folder.

```sh
go get github.com/Moonlington/harmonia
```

This pulls the latest tagged release from the master branch.

### Usage

Import **Harmonia** into your project.

```go
import "github.com/Moonlington/harmonia"
```

You can now create a new **Harmonia** client, which holds most of the features of this project in addition to having the same functionality of a [DiscordGo session](https://pkg.go.dev/github.com/bwmarrin/discordgo#Session).

```go
harmonia, err := harmonia.New("authentication token")
```

Please refer to the [documentation](https://pkg.go.dev/github.com/Moonlington/harmonia) and take a look at the [examples](https://github.com/Moonlington/harmonia/tree/main/examples) for more detailed information.
