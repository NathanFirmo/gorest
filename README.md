<div align=center>

# Gorest

![Open Source Love](https://badges.frapsoft.com/os/mit/mit.svg?v=102)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)

Enjoy making HTTP requests in your terminal, just like you do in Insomnia.

<img src="https://github.com/NathanFirmo/gorest/assets/79997705/1cd5f526-9e34-4493-810b-f0f40b369986" width="550px" heigth="550px">

</div>

## Demo video

https://github.com/NathanFirmo/gorest/assets/79997705/85f7bf5d-2f23-41a1-b605-9bb2fee441b1

## Installing

Binaries available on the [releases page](https://github.com/NathanFirmo/gorest/releases/latest).

Or, with Go installed:

```sh
go install github.com/NathanFirmo/gorest/cmd/gorest@latest
```

## Commands

|Key|Action|
|---|---|
|`CTRL-H`|Go to section on the **left**|
|`CTRL-J`|Go to field **below**|
|`CTRL-K`|Go to field **above**|
|`CTRL-L`|Go to section on the **right**|
|`CTRL-N`|Creates a **new request** on the list|
|`CTRL-Space` / `CTRL-R`|**Execute** current request|
|`CTRL-O`|Toggle response **body/headers**|
|`CTRL-T`|Switch to the next **theme**|
|`d` / `Delete`|**Delete** the selected request (on the list)|
|`CTRL-C`|Quit|

## Importing curl commands

Paste a curl command into the **URL** field, or type one and press `Enter`, and gorest fills in the method, URL, headers and body of the current request. Commands copied from the browser DevTools ("Copy as cURL") work too.

## Themes

`dracula` (default), `catppuccin`, `tokyo-night`, `gruvbox`, `nord` and `solarized-light`. Press `CTRL-T` to cycle through them; your choice is remembered.

## Data

Requests are stored in SQLite at `$XDG_DATA_HOME/gorest/gorest.db` (usually `~/.local/share/gorest/gorest.db`; the user config directory on macOS and Windows). Set `GOREST_DB` to use another file. Databases created by older versions at `~/gorest/.db` are moved there automatically.

## Roadmap

- [x] ~Save requests locally with SQLite~;
- [x] ~Allow pasting curl requests like Insomnia does~;
- [x] ~Add themes~;
- [ ] Support multipart forms (`curl -F`);
- [ ] Environments and variables;
- [ ] Add gorest to some package managers like flatpak, snap, aur, winget, etc.
