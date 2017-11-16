# alfred-project-manager

Allows you to quickly open projects from Alfred.
Set `PROJECT_DIRECTORIES` inside Alfred's environment variables, and you're good to go.

### Usage

```
pm <fuzzy search term>
```

### Modifiers

* `none` (default): open the project in VSCode
* `alt`: open the project in Finder
* `cmd`: open the project in iTerm

You can change which applications are used to open the project inside the workflow.
