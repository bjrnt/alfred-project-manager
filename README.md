# alfred-project-manager

> Allows you to quickly open projects from Alfred.

![usage example](/image.png)

## Installation

- Grab the latest release [here](https://github.com/bjrnt/alfred-project-manager/releases/) and install the workflow file.

## Usage

Open Alfred and type `pm` to access the project manager and try typing a query. You can also configure a hotkey for it by opening the workflow in Alfred's workflow options panel.

### Modifiers

- `none` (default): open in editor (default: VSCode)
- `alt/opt`: reveal in Finder
- `cmd`: open in terminal (default: iTerm)
- `ctrl`: open the project's repo in your browser

You can change the modifier and application combinations in Alfred's workflow settings window.

### Features

- Enable "Show only valid Git repos" if you want to exclude folders in your project directory that are not Git repos.
- Increase "Maximum Project Depth" if your projects are grouped within other folders in your projects directory. For example, if my "projects" folder has folders "personal" and "work" and both of these contain my projects, I'd set the value to 1.

## Maintenance

### Issues and Feature Requests

Feel free to open an issue for the project if you have encountered a problem or have a feature request for the workflow.

### Building

The project can be built, linked to Alfred, and released using [jason0x34/go-alfred](https://github.com/jason0x43/go-alfred). Commands for this can be found in [./.vscode/tasks.json](./.vscode/tasks.json).
