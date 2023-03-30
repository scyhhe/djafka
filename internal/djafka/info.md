# Info

1. Select your preferred connection from the connection pane
2. Select the entity you want to know more about
3. Look at the information
4. Success :)

# Overview

Support for multiple Kafka Cluster connections via config file.

These are specified in the ```config.json``` in the base directory.

```json
{
    "connections": [
        {
            "name": "localhost",
            "bootstrapServer": "localhost"
        }
    ]
}
```

# Help and Navigation

```?``` Toggle Help context menu with available shortcuts

```q``` Quits the CLI tool

You can use both arrows and VIM shortcuts for navigation.

```←/h``` move left

```↓/j``` move down

```↑/k``` move up

```→/l``` move right

You can cycle between windows/panes with `TAB`

```ctrl + t``` create a new topic in the currently active connection

```ctrl + r``` reset the offset for the currently selected topic
