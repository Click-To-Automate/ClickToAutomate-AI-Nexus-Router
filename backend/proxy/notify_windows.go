package proxy

import (
	"fmt"
	"log"
	"os/exec"
)

// TriggerDesktopNotification sends a non-intrusive Windows toast notification
// using a background PowerShell process.
func TriggerDesktopNotification(failedProvider string, fallbackProvider string, statusCode int) {
	msg := fmt.Sprintf("%s failed (%d). Switched to %s.", failedProvider, statusCode, fallbackProvider)
	
	psScript := fmt.Sprintf(`
[reflection.assembly]::loadwithpartialname("System.Windows.Forms") | Out-Null
$notify = new-object system.windows.forms.notifyicon
$notify.icon = [System.Drawing.SystemIcons]::Warning
$notify.balloonTipIcon = "Warning"
$notify.balloonTipTitle = "CTA AI Router Fallback"
$notify.balloonTipText = "%s"
$notify.visible = $true
$notify.showBalloonTip(5000)
Start-Sleep -Seconds 5
$notify.visible = $false
$notify.Dispose()
`, msg)

	cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", psScript)
	err := cmd.Start()
	if err != nil {
		log.Printf("[Notifier] Failed to trigger desktop notification: %v\n", err)
	}
}
