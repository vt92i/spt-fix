package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/godbus/dbus/v5"
)

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	if err = conn.AddMatchSignal(
		dbus.WithMatchSender("org.mpris.MediaPlayer2.spotify"),
		dbus.WithMatchObjectPath("/org/mpris/MediaPlayer2"),
		dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		dbus.WithMatchMember("PropertiesChanged"),
		dbus.WithMatchOption("arg0", "org.mpris.MediaPlayer2.Player"),
	); err != nil {
		panic(err)
	}

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
	for v := range c {
		md, ok := v.Body[1].(map[string]dbus.Variant)["Metadata"].Value().(map[string]dbus.Variant)

		if !ok {
			continue
		}

		artist := md["xesam:artist"].Value().([]string)[0]
		title := md["xesam:title"].Value()

		current_song := fmt.Sprintf("%s - %s", artist, title)
		fmt.Println(current_song)

		set_window_title(current_song)
	}

}

func set_window_title(title string) {
	// xdotool search --class "spotify" set_window --name {{title}}
	program := "/usr/bin/xdotool"
	arg1 := "search"
	arg2 := "--class"
	arg3 := "spotify"
	arg4 := "set_window"
	arg5 := "--name"
	arg6 := title

	exec.Command(program, arg1, arg2, arg3, arg4, arg5, arg6).Run()
}
