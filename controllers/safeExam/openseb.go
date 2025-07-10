package controllers

import (
	"bytes"
	"fmt"
	"os/exec"

	"wakeonlan/config"

	"github.com/gin-gonic/gin"
)


type PC struct {
	Name         string
	IPAddress    string
	ComputerName string
	Username     string
	Password     string
}


func LaunchSEBHandler(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT ip_address, computer_name, username, password
		FROM pcs
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": "DB query failed: " + err.Error()})
		return
	}
	defer rows.Close()

	var pcs []PC
	for rows.Next() {
		var pc PC
		if err := rows.Scan(&pc.IPAddress, &pc.ComputerName, &pc.Username, &pc.Password); err != nil {
			c.JSON(500, gin.H{"error": "DB scan failed: " + err.Error()})
			return
		}
		pcs = append(pcs, pc)
	}

	results := make(map[string]string)

	for _, pc := range pcs {
		
		psScript := fmt.Sprintf(`Invoke-Command -ComputerName "%s" -ScriptBlock {
    $Action = New-ScheduledTaskAction -Execute "C:\Program Files\SafeExamBrowser\Application\SafeExamBrowser.exe"
    $Trigger = New-ScheduledTaskTrigger -AtLogOn
    $Principal = New-ScheduledTaskPrincipal -UserId "%s" -LogonType Interactive -RunLevel Highest
    Register-ScheduledTask -TaskName "LaunchSEB" -Action $Action -Trigger $Trigger -Principal $Principal -Force
    Start-ScheduledTask -TaskName "LaunchSEB"
} -Credential (New-Object System.Management.Automation.PSCredential("%s", (ConvertTo-SecureString "%s" -AsPlainText -Force)))`, pc.IPAddress, pc.Username, pc.Username, pc.Password)

		cmd := exec.Command("powershell", "-Command", psScript)
		var out, stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			results[pc.IPAddress] = fmt.Sprintf("PowerShell error: %s", stderr.String())
		} else {
			results[pc.IPAddress] = "SEB Launch success"
		}
	}

	c.JSON(200, gin.H{"results": results})
}
