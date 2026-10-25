package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

const usage = `thctl - send commands to the thunder music player

Usage:
  echo COMMAND | thctl [-p path]
  thctl -h

Options:
  -h         Print this help and exit.
  -p path    Path to the unix socket (default: /tmp/thunder.sock).

Commands:
  TOGGLE_PAUSE   Pause or resume playback.
  TOGGLE_MUTE    Mute or unmute the sound.
  PREVIOUS       Play the previous track.
  NEXT           Play the next track.
  VOL_DOWN       Decrease the volume.
  VOL_UP         Increase the volume.
  STOP           Stop playback.
  CYCLE_MODE     Switch to the next playback mode.

Examples:
  echo "NEXT" | thctl
  echo "VOL_UP" | thctl -p /tmp/thunder.sock
`

var (
	path = flag.String("p", "/tmp/thunder.sock", "")
	help = flag.Bool("h", false, "")
)

func init() {
	flag.Usage = func() {
		fmt.Printf("%s\n", usage)
	}
	flag.Parse()

	if *help == true {
		flag.Usage()
		os.Exit(0)
	}
}

type Response struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Length   string `json:"length"`
	Volume   int    `json:"volume"`
	Muted    bool   `json:"muted"`
	Playing  bool   `json:"playing"`
	Paused   bool   `json:"paused"`
}

func main() {
	if err := write(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}

func write() error {
	conn, err := net.Dial("unix", *path)
	if err != nil {
		return err
	}
	defer conn.Close()

	buff := make([]byte, 64)
	n, err := os.Stdin.Read(buff)
	if n > 0 {
		if _, err := conn.Write(buff[:n]); err != nil {
			return fmt.Errorf("write: %s\n", err)
		}
	}
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
