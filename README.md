# alfred-project-manager

![usage example](/image.png)

Allows you to quickly open projects from Alfred. It will try its best to provide the best possible match for your query.

## Installation

Grab the latest release [here](https://github.com/bjrnt/alfred-project-manager/releases/) and install the workflow file.

Set `PROJECT_DIRECTORY` inside the workflow's environment variables (install the workflow, select in from the workflow list, and press the `[x]` icon), and you're good to go.

## Usage

Open Alfred and type `pm` to access the project manager and try typing a query. You can also configure a hotkey for it.

## Modifiers

- `none` (default): open the project in VSCode
- `alt`: open the project in Finder
- `cmd`: open the project in iTerm
- `ctrl`: open the project's repo

You can change which applications are used to open the project inside the workflow.
