# alfred-project-manager

> Allows you to quickly open projects from Alfred.

![usage example](/image.png)

## Installation

1. Grab the latest release [here](https://github.com/bjrnt/alfred-project-manager/releases/) and install the workflow file.
2. Set the `PROJECT_DIRECTORY` path inside the workflow's environment variables (select the workflow and press the `[x]` icon in the upper right). A relative path is assumed to be relative to your home directory.

## Usage

Open Alfred and type `pm` to access the project manager and try typing a query. You can also configure a hotkey for it by opening the workflow in Alfred's workflow options panel.

### Modifiers

- `none` (default): open in editor (default: VSCode)
- `alt/opt`: reveal in Finder
- `cmd`: open in terminal (default: iTerm)
- `ctrl`: open the project's repo in your browser

You can change the modifier and application combinations in Alfred's workflow settings window.

### Further Config

#### Show Only Git Repos

If you only want the workflow to show valid Git repos, you can set the `REQUIRE_DOTGIT` to `true` in the workflow's environment variables.

#### Workspaced Projects

You can configure the workflow `MAX_PROJECT_DEPTH` to an integer larger than zero if you have nested or workspaced projects. For example, if you have different folders for different companies or clients where you have your projects, you would set this to 1.

It is recommended you combine this option with only showing git repos to avoid false positives.

## Maintenance

### Issues and Feature Requests

Feel free to open an issue for the project if you have encountered a problem or have a feature request for the workflow.

### Building

The project can be built, linked to Alfred, and released using [jason0x34/go-alfred](https://github.com/jason0x43/go-alfred). Commands for this can be found in [./.vscode/tasks.json](./.vscode/tasks.json).
