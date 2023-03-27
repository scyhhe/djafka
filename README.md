# djafka

## Getting Started

To get started you need Go >= 1.20 installed.

Run `make init` to initialize the project and downloads all dependencies.
Running `make` builds the final `djafka` executable which can be found under
`bin/djafka`.

### Overview

Create an extremely easy-to-use, intuitive, interactive, and keyboard friendly CLI Tool to interact with a Kafka Cluster.

Inspired by Philipp and Lazygit.

### Goals

- Connect to a kafka server
- List topics
- List consumers
- Add new topics
- Show consumer offsets
- Reset consumer offsets  
- Create a prototype with these features
- Create a presentation about our findings and learnings

### Non Goals

- Make it connect to MSK in our environment
- Make it run on a bastion host

#### Optional

- Probe topics
- Interactively show topic events consumption by consumers
- Show cluster stats
- Display and collect metrics
- Make polling timeframes configurable
- Notify when a consumer caught up to the latest offset

### Proposed Solution

In the spirit of the Hackathon, we decided to go with ... Go! 🫠  as our language of choice.

#### CLI Libraries

There are a number of Go! TUI tools, each with their own pros and cons. In general, what we want to support is

- Frames / Panes
- Selectable Lists
- Shortcut Support (Mouse Support would be nice also)
- Coloring
- Graphs / Charts / Progress Bars / Gauges
- User Input and Prompts

We evaluted the following libraries

- [termdash](https://github.com/mum4k/termdash)
- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [termui](https://github.com/gizak/termui)
- [gocui](https://github.com/jroimartin/gocui)

| Library   | Lists | Shortcuts           | Colors | Graphs   | User Input | Frames | Notes                  |
| --------- | ----- | ------------------- | ------ | -------- | ---------- | ------ | ---------------------- |
| termdash  | ❌     | ✅                   | ✅      | ✅        | ✅          | ✅      | Not as lit             |
| bubbletea | ✅     | ✅ (+ Mouse Support) | ✅      | ✅        | ✅          | ✅      | Lit                    |
| termui    | ✅     | ✅                   | ✅      | ✅(a lot) | ✅          | ✅      | Outdated and kinda lit |
| gocui     | ✅     | ✅                   | ✅      | ❌        | ✅          | ✅      | not lit                |

In the end, we decided to go with [bubbletea](https://github.com/charmbracelet/bubbletea) because of its superior features set, such as a defined way of doing state updates, a lot of functional examples, logging, and good documentation.

#### Kafka Admin Tool

For administrating Kafka, we decided to use the official [Confluent Kafka Go Client](https://github.com/confluentinc/confluent-kafka-go), since it provides almost everything we need right out of the box.

Docs can be found [here](https://docs.confluent.io/platform/current/clients/confluent-kafka-go/index.html).
