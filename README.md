# tx
Minimalistic tasks for CLI - with useful extras!

`tx` stands for *t extended*, and is a complete rewrite of [Steve Losh's *t*](https://stevelosh.com/projects/t/). `tx` tries to retain the original project's minimalism and its core principle:

> [...]the only way to make your todo list prettier is to **finish some damn tasks**.

...while introducing new features that make it more pleasant *(but admittedly, more bloated)* to work with in a command line interface.

# Features

- ✔️ **Written in Go**: `tx` is a complete reimplementation of `t`. While the original projects was written in Python, `tx` uses the Go programming language, which makes it **faster and more portable** (since it's a single compiled binary).
- ⏪ **Backwards Compatible**: Even though `tx` is a complete rewrite, file handling and formatting is the exact same. `tx` generates more metadata (like creation date and time), but also includes a *"legacy"* `id` attribute - meaning **tasklists generated by `t` can be used with `tx`** and vice versa.
- 🧠 **Multiple Actions**: Gone are the days of invoking multiple instances of `t` when you want to remove a task, finish a task then add a new one in this sequence. **`tx` handles action flags as sequencial commands**, meaning you can accomplish the above example this way: `tx --remove 1 --finish 2 --add "Easy!"`. You can also specify the **same action multiple times**!
- 📝 **Smarter Task Selection**: While you *can* specify multiple `--remove` commands to remove multiple tasks, `tx`'s "selector" system allows for more flexibility, allowing the selection of mulitple tasks in a variety of formats. **Convenience actions** like `--wipe` and `--complete` are also present to erase a full tasklist or mark all tasks as finished. Further information in [Selectors](#selectors).
- 🌐 **Syncing**: `tx` offers a simple way to **move tasklists between devices and keep them in sync**. By default, the excellent [JSON Blob](https://jsonblob.com/) is used as the backbone of syncing, but you can **run your own syncing service** by mimicking the JSON Blob API and changing the Sync service using `tx sync change`. More information in [Syncing Details](#syncing-details).
- 🤝 **Script Friendly**: `tx` was designed to be useful in a software development/CLI environment, and the use of `tx` is encouraged in shell scripts by featuring **the specification of custom output formats, the running of a callback shell command when the tasklist is modified through `tx` and many documented exit codes**. See [Scripting Help](#scripting-help) for more information.

# Table of Contents

1. [Setup](#setup)
2. [General Usage](#general-usage)
3. [Syncing Details](#syncing-details)
4. [Scripting Help](#scripting-help)
5. [Contributions](#contributions)

# Setup

1. **Download the latest binary from the [Releases Page](#https://github.com/doczi-dominik/tx/releases) or clone the repository and build it yourself with Go:**

```sh
$ git clone https://github.com/doczi-dominik/tx
[...]
$ cd tx
$ go get
$ go install  # Use this if you have GOBIN correctly configured
$ go build    # Use this if you just want to compile a binary.
```

2. (optional) **Add aliases for tx to *"set default options"***. For example, if you use Bash, type this in `.bashrc`:

```sh
# Tasks mode
alias t='tx --list="tasks" tasks'

# Finished tasks mode
alias td='tx --list="tasks" --delete-if-empty done'
```

# General Usage

*Note: The usage section uses the aliases defined in [Setup](#setup)/2. Some of the actions are only applicable to tasks mode or done mode. Check the help page with `--help/-h` to learn more.*

## Modes

`tx` uses subcommands, just like Git! These are:
- `tx tasks`: List and modify active tasks
- `tx done`: List and modify finished tasks
- `tx sync`: Configure syncing for the current tasklist

Pass `--help/-h` after passing the mode (or take a look at the [Wiki](https://github.com/doczi-dominik/tx/wiki)) to learn more.

## Listing Tasks

To list tasks, just use `tx tasks` or `tx done` without any extra arguments:

```sh
$ t
1 - these
2 - are
3 - active
4 - tasks

$ td
1 - these (2021/01/14)
2 - are (2021/01/14)
3 - finished (2021/01/14)
4 - tasks (2021/01/14)
```

## Adding Tasks

To add tasks in `tasks` mode, either:
- Enter text after the last argument
- Use the dedicated `--add/-a` action

```
$ t You do not need quotes here!
1 - You do not need quotes here!

$ t --add "Use this when you pass other arguments after this." --add "Also, you need quotes here for multiple words!"
1 - You do not need quotes here!
2 - Use this when you pass other arguments after this.
3 - Also, you need quotes here for multiple words!
```

## Selectors

You can use selectors with any of `--remove`, `--finish`, `--edit` in `tasks` mode and `--restore` and `--delete` in `done` mode.

Tasks can be selected using their **indexes**. A task's index is **fixed to a task**, meaning that indexes don't shift around when removing and finishing tasks. Internally, indexes are only recalculated before writing them into the taskfile and displaying them.

Tasks can be selected using their indexes, printed before their names (by default) or by using **placeholder letters** instead of using numbers. The letters `f` and `r` always refer to the **first task**, and `l` always refers to the **last task.**

The following formats can be used:
- **Index**: A single task's index. Examples:
    - 2
    - f
    - R
- **CSV**: Comma-separated indexes. Examples:
    - 3,4,5
    - f,7,L
- **Range:** A range of indexes. Examples:
    - 1-6
    - 4-2 *(Note: The range is reversed by `tx` to 2-4.)*
    - f-l *(Note: You can also use convenience actions like `--wipe` or `--complete` for operations on every task.)*

Please note that when using `--edit`, **only Index notation can be used** as editing multiple tasks may make duplicates or cause other issues.

## Editing a task

You can edit tasks in two ways, either replacing the full task text (equivalent to removing and adding a new task, but keeping its index) or by replacing all occurences of a word in it:

```
$ t
1 - I am a task.
2 - hello hello hello

$ t --edit "f/First Task" --edit "l/hello/Howdy"
1 - First Task
2 - Howdy Howdy Howdy
```

## Enabling Syncing

To enable syncing for a particular tasklist, use `tx sync enable`. By default, this will request a new, unique Sync ID from the default Sync service. To connect your tasklist with an existing Sync ID, write it after the command like so: `tx sync enable "this-is-the-sync-id"`. Read the [Wiki](https://github.com/doczi-dominik/tx/wiki) for details on the `sync` mode.

# Syncing Details

`tx`'s syncing feature uses JSON and a simple HTTP server as its core. `tx` is designed to use Tristan Burch's [JSON Blob API](https://jsonblob.com/api), but any HTTP server can be provided which uses the following model:

## Tasklist Format

Tasklists are accessed with their respective Sync ID's from the server. Taskfiles are serialized into a simple JSON object with two keys: `contents` and `doneContents`. Both keys point to a single string, which contains the respective active tasklist and finished tasklists.

### Example

- `tasks.txt`:
```
one taskline
two tasklines
three tasklines

```

- `.tasks.txt.done`:
```
one finished taskline
two finished tasklines

```

- Serialzed JSON Object:
```json
{
    "contents": "one taskline\ntwo tasklines\nthree tasklines\n",
    "doneContents": "one finished taskline\ntwo finished tasklines\n"
}
```

## Server API

As mentioned above, the server API should mimic the [JSON Blob API](https://jsonblob.com/api). The Sync ID of the tasklist will be appended to the provided Sync URL (e.g.: `"https://jsonblob.com/api/jsonBlob/" + syncID`) and `tx` will use the appropriate HTTP request depending on the operation.

## Sync Information Storage

The necessary information used for syncing is stored in a *syncfile*, the filename is the taskfile's filename, but a `.` is prepended to ensure it's a hidden a file, and `.sync` is appended to signify that this is a syncfile.

Syncfiles contain the following information:
- `syncID`: The Sync ID of the tasklist as requested from the sync service
- `syncURL`: If present, it will override the `--sync-url/-U` flag.
- `lastNetworkUpdate`: A date signifying the time of the last successful POST/PUT request. This will be compared with a local taskfile's `mtime` to determine if the synced tasklist is up-to-date.

# Scripting help

## Callback

You can pass a command line to `--callback/-C` just like you would type it in a shell. The callback will only run if the tasklist was modified. The provided command's STDOUT and STDERR streams are combined into `tx`'s STDOUT stream. For more information, read [Go's 'exec' module documentation on cmd.CombinedOutput()](https://pkg.go.dev/os/exec#Cmd.CombinedOutput).

## Output Formatting

The listing output can be customized using the `--output/-o` flag. The following placeholders are available:
- `{index}`: The index of a given task
- `{task}`: The task text
- `{creationDate}`: The date of the task's creation in `YYYY/MM/DD` format.
- `{creationTime}`: The time of the task's creation in `HH:MM` format.
- `{finishedDate}`: The date the task was marked as finished in `YYYY/MM/DD` format.
- `{finishedTime}`: The time the task was marked as finished in `HH:MM` format.

If you want to use any of these placeholders *literally*, simply double the braces: `{{index}}`. This will show up as `{index}` in the output, not as the task's index.

## Exit Codes

### Miscellaneous

Code | Meaning
---- | -------
0 | Success, no errors
1 | Command line flag parsing failed
2 | Invalid selector notation
3 | Invalid selector notation for `--edit/-e`
4 | Invalid index, no task exists at that index
5 | Task validation failed (e.g.: task contains a newline in the middle)
6 | Non-creative operation on empty tasklist
7 | Cannot serialize to JSON
8 | Callback could not execute / returned a non-zero error code

### Taskfile Operations

Code | Meaning
---- | -------
9 | Could not open taskfile
10 | Could not write taskfile
11 | Could not read taskfile

### Backup Taskfile Operations

Code | Meaning
---- | -------
12 | Could not create backup file
13 | Could not write backup file

### Syncfile Operations

Code | Meaning
---- | -------
14 | Could not open syncfile
15 | Could not write syncfile
16 | Could not read syncfile

### HTTP GET Requests

Code | Meaning
---- | -------
17 | Could not create GET Request

### HTTP PUT Requests

Code | Meaning
---- | -------
18 | Could not create PUT Request
19 | Could not complete PUT Request

### HTTP DELETE Requests

Code | Meaning
---- | -------
20 | Could not create DELETE Request
21 | Could not complete DELETE Request

### Sync Server Errors

Code | Meaning
---- | -------
22 | No tasklist exists with this Sync ID
23 | Could not request new Sync ID from the Sync server
24 | Unparseable response from Sync server
25 | Unsupported Configuration, mainly exists to signify that deleting this blob is disabled on the Sync service, which `tx` will never configure.

# Contributions

Issues and PRs are always welcome, be it as small as a typo or as large as a new feature!
