# alfred-project-manager
> Allows you to quickly open projects from Alfred.

![usage example](/image.png)

## Installation

1. Grab the latest release [here](https://github.com/bjrnt/alfred-project-manager/releases/) and install the workflow file.
2. Set the `PROJECT_DIRECTORY` path inside the workflow's environment variables (select the workflow and press the `[x]` icon in the upper right). A relative path is assumed to be relative to your home directory.

## Usage

Open Alfred and type `pm` to access the project manager and try typing a query. You can also configure a hotkey for it by opening the workflow in Alfred's workflow options panel.

## Modifiers

- `none` (default): open the project in VSCode
- `alt/opt`: open the project in Finder
- `cmd`: open the project in iTerm
- `ctrl`: open the project's repo

You can change which applications are used to open the project inside the workflow.
