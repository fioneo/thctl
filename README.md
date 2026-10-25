# thctl

thctl is a command-line client for controlling the thunder player. It reads a command from standard input and sends it to the player over a unix socket.

## Build and install

From the directory containing the thctl source:

```sh
go build -o thctl .
install -m 755 thctl ~/.local/bin/thctl
```

Make sure the install directory is in `PATH`, since window manager key bindings call `thctl` by name.

## Usage

```
echo COMMAND | thctl [-p path]
thctl -h
```

Options:

| Option    | Description                                              |
|-----------|----------------------------------------------------------|
| `-h`      | Print help and exit.                                     |
| `-p path` | Path to the unix socket. Default: `/tmp/thunder.sock`.   |

Example:

```sh
echo "NEXT" | thctl
echo "VOL_UP" | thctl -p /tmp/thunder.sock
```

## Commands

| Command        | Action                              |
|----------------|-------------------------------------|
| `TOGGLE_PAUSE` | Pause or resume playback.           |
| `TOGGLE_MUTE`  | Mute or unmute the sound.           |
| `PREVIOUS`     | Play the previous track.            |
| `NEXT`         | Play the next track.                |
| `VOL_DOWN`     | Decrease the volume.                |
| `VOL_UP`       | Increase the volume.                |
| `STOP`         | Stop playback.                      |
| `CYCLE_MODE`   | Switch to the next playback mode.   |

Commands are case-sensitive. A trailing newline is ignored.

## Window manager integration

### i3

Add to `~/.config/i3/config`:

```
bindsym $mod+F1 exec --no-startup-id echo "TOGGLE_PAUSE" | thctl
bindsym $mod+F2 exec --no-startup-id echo "TOGGLE_MUTE" | thctl
bindsym $mod+F3 exec --no-startup-id echo "PREVIOUS" | thctl
bindsym $mod+F4 exec --no-startup-id echo "NEXT" | thctl
bindsym $mod+F5 exec --no-startup-id echo "VOL_DOWN" | thctl
bindsym $mod+F6 exec --no-startup-id echo "VOL_UP" | thctl
bindsym $mod+F7 exec --no-startup-id echo "STOP" | thctl
bindsym $mod+F8 exec --no-startup-id echo "CYCLE_MODE" | thctl
```

Reload with `$mod+Shift+r`.

### dwm

Add to `config.h`. The macro goes near the top, next to the other helper macros:

```c
#define THCTL(cmd) SHCMD("echo " cmd " | thctl")
```

The key bindings go into the `keys[]` array:

```c
{ MODKEY, XK_F1, spawn, THCTL("TOGGLE_PAUSE") },
{ MODKEY, XK_F2, spawn, THCTL("TOGGLE_MUTE") },
{ MODKEY, XK_F3, spawn, THCTL("PREVIOUS") },
{ MODKEY, XK_F4, spawn, THCTL("NEXT") },
{ MODKEY, XK_F5, spawn, THCTL("VOL_DOWN") },
{ MODKEY, XK_F6, spawn, THCTL("VOL_UP") },
{ MODKEY, XK_F7, spawn, THCTL("STOP") },
{ MODKEY, XK_F8, spawn, THCTL("CYCLE_MODE") },
```

Rebuild and restart dwm:

```sh
sudo make clean install
```

If `SHCMD` is not defined in your `config.h`, add it:

```c
#define SHCMD(cmd) { .v = (const char*[]){ "/bin/sh", "-c", cmd, NULL } }
```

### Other tools

Any tool that can run a shell command works the same way, for example sxhkd:

```
super + F4
    echo "NEXT" | thctl
```

## Troubleshooting

- `error: dial unix /tmp/thunder.sock: connect: no such file or directory` - the player is not running, or it uses a different socket path. Pass the correct path with `-p`.
- `error: dial unix ...: connect: connection refused` - a stale socket file is left over from a previous run. Remove it and restart the player.
- The command is sent but nothing happens - check the player log. Unknown commands are ignored by the server.
