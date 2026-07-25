# Keybindings

Volt's terminal interface uses contextual, Vim-inspired keybindings. Press `?`
outside a text editor to see shortcuts for the focused panel, or press `F1`
anywhere to see every shortcut.

## Key notation

- `Ctrl+key` means hold Control while pressing the key.
- `Alt+key` means Alt on Linux and Windows, or Option on macOS.
- `Esc Esc` means press Escape twice within 750 milliseconds.
- Uppercase bindings such as `G` require Shift.

## Global

| Keys | Action |
|---|---|
| `Ctrl+C` | Quit immediately and cancel an active load test |
| `Esc Esc` | Quit after a deliberate double press |
| `F1` | Show all shortcuts |
| `Alt+H` | Focus the previous panel |
| `Alt+L` | Focus the next panel |

The first `Esc` closes help or returns focus to the sidebar and displays a quit
prompt. A different key or the 750-millisecond timeout cancels the quit
sequence.

`q` never quits Volt. It is inserted normally while editing text, does nothing
in non-editing panels, and closes the help overlay when the overlay is open.

### Alt on macOS

macOS calls the Alt key `Option`. The panel bindings are therefore
`Option+H` and `Option+L`.

Your terminal must send Option as Alt/Meta, sometimes described as sending an
Escape prefix. If `Option+H` produces a special character such as `˙`, enable
the terminal profile's Option-as-Meta or `Esc+` setting. This is a terminal
configuration; Volt receives the same `Alt+H` event on every platform.

## Sidebar

| Keys | Action |
|---|---|
| `Enter` | Open the selected saved request |
| `j` | Select the next request |
| `k` | Select the previous request |
| `g` | Select the first request |
| `G` | Select the last request |
| `d` | Delete the selected request |
| `?` | Show sidebar shortcuts |

## Request panel

| Keys | Action |
|---|---|
| `Tab` | Focus the next field |
| `Shift+Tab` | Focus the previous field |
| `h` | Select the previous HTTP method |
| `l` | Select the next HTTP method |
| `Ctrl+Enter` or `Alt+Enter` | Send the request |
| `Ctrl+S` | Save the request |
| `Ctrl+L` | Toggle load-test mode |
| `?` | Show request shortcuts when not editing |

Plain `Enter` remains available to edit content or accept suggestions. Arrow
keys remain available to text editors. Request bodies are edited as raw text,
so JSON, XML, GraphQL, and other formats do not require key/value conversion.

While a text field is focused, printable shortcut characters such as `?`, `q`,
`h`, and `l` are inserted into the field instead of triggering panel actions.

## Response panel

| Keys | Action |
|---|---|
| `h` | Select the previous response tab |
| `l` | Select the next response tab |
| `1`–`3` | Jump directly to a response tab |
| `j` | Scroll down |
| `k` | Scroll up |
| `Ctrl+D` | Scroll down half a page |
| `Ctrl+U` | Scroll up half a page |
| `y` | Copy the response body |
| `Ctrl+X` | Cancel an active load test |
| `?` | Show response shortcuts |

## Help overlay

| Keys | Action |
|---|---|
| `h` | Show the previous help section |
| `l` | Show the next help section |
| `1`–`4` | Jump directly to a help section |
| `q` or `?` | Close help |
| `Esc` | Close help and arm the double-Escape quit sequence |

The in-application help is generated from the same action registry used for
input dispatch. When a binding changes in code, its displayed help and
contextual binding tests change with it.
