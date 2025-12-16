package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed icon.ico
var iconData []byte

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("")
	systray.SetTooltip("FaveoAgent Controller")

	// Menu items
	restart := systray.AddMenuItem("Restart FaveoAgent", "Restart the FaveoAgent service")
	openDash := systray.AddMenuItem("Open Dashboard", "Open Faveo Dashboard in browser")
	createTicket := systray.AddMenuItem("Create Ticket", "Open ticket creation UI")
	stop := systray.AddMenuItem("Stop FaveoAgent", "Stop the FaveoAgent service")
	uninstall := systray.AddMenuItem("Uninstall FaveoAgent", "Completely uninstall FaveoAgent")
	quit := systray.AddMenuItem("Quit", "Quit the app")

	// Handle click events
	go func() {
		for {
			select {
			case <-restart.ClickedCh:
				runCommand("systemctl", "restart", "faveoagent.service")

			case <-openDash.ClickedCh:
				openBrowser("https://google.com")

			case <-createTicket.ClickedCh:
				openBrowser(fmt.Sprintf("https://agentsw.faveodemo.com/create-ticket?agent_id=%s", os.Getenv("FAVEO_AGENT_ID")))

			case <-stop.ClickedCh:
				runCommand("systemctl", "stop", "faveoagent.service")

			case <-uninstall.ClickedCh:
				runCommand("/usr/local/bin/faveoagent")

			case <-quit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\nOutput: %s\n", err, string(output))
	} else {
		fmt.Printf("Command output: %s\n", string(output))
	}
}

// openExternalBrowser opens URL in system default browser
func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		bashCmd := fmt.Sprintf(`
user=$(who | awk '/tty|pts/ {print $1; exit}')
uid=$(id -u "$user")
export XDG_RUNTIME_DIR="/run/user/$uid"
export DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/$uid/bus"
sudo -u "$user" xdg-open "%s"
`, url)

		cmd = exec.Command("bash", "-c", bashCmd)

	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)

	case "darwin":
		cmd = exec.Command("open", url)

	default:
		fmt.Println("Unsupported platform")
		return
	}
	cmd.Start()
}
