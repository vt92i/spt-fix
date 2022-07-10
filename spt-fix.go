package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func main() {
	// playerctl --player=spotify metadata --format '{{title}} - {{artist}}' --follow
	program := "/usr/bin/playerctl"
	arg1 := "--player=spotify"
	arg2 := "metadata"
	arg3 := "--format"
	arg4 := "'{{title}} - {{artist}}'"
	arg5 := "--follow"

	cmd := exec.Command(program, arg1, arg2, arg3, arg4, arg5)
	stdout, _ := cmd.StdoutPipe()

	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		current_song := scanner.Text()

		if len(current_song) == 0 {
			continue
		}

		if len(current_song) == 5 {
			for {
				out, _ := exec.Command(program, arg1, arg2, arg3, arg4).Output()
				current_song = string(out)
				if len(current_song) > 6 {
					current_song = current_song[1 : len(current_song)-2]
					break
				}
			}
		} else {
			current_song = current_song[1 : len(current_song)-1]
		}

		set_window_title(current_song)
		fmt.Println(current_song)
	}

	cmd.Wait()
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
