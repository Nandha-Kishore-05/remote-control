// package controllers

// import (
// 	"fmt"
// 	"net/http"
// 	"os/exec"
// 	"wakeonlan/config"

// 	"github.com/gin-gonic/gin"
// )

// type PCs struct {
// 	Name         string
// 	IPAddress    string
// 	ComputerName string
// 	Username     string
// 	Password     string
// }

// func ExitSEBHandler(c *gin.Context) {
// 	rows, err := config.DB.Query(`
// 		SELECT name, ip_address, computer_name, username, password
// 		FROM pcs
// 	`)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed: " + err.Error()})
// 		return
// 	}
// 	defer rows.Close()

// 	var pcs []PC
// 	for rows.Next() {
// 		var pc PC
// 		if err := rows.Scan(&pc.Name, &pc.IPAddress, &pc.ComputerName, &pc.Username, &pc.Password); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB scan failed: " + err.Error()})
// 			return
// 		}
// 		pcs = append(pcs, pc)
// 	}

// 	results := make(map[string]interface{})

// 	for _, pc := range pcs {
// 		// Build PowerShell script to force kill and remove RecoveryState.seb
// 		forceKillScript := fmt.Sprintf(`
// $secpasswd = ConvertTo-SecureString "%s" -AsPlainText -Force
// $cred = New-Object System.Management.Automation.PSCredential ("%s", $secpasswd)
// Invoke-Command -ComputerName "%s" -Credential $cred -Authentication Credssp -ScriptBlock {
//     taskkill /im SafeExamBrowser.exe /f
//     Remove-Item "C:\ProgramData\SafeExamBrowser\RecoveryState.seb" -ErrorAction SilentlyContinue
// }
// `, pc.Password, pc.Username, pc.IPAddress)

// 		// Execute the force kill script
// 		cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", forceKillScript)
// 		output, err := cmd.CombinedOutput()

// 		if err != nil {
// 			results[pc.Name] = gin.H{
// 				"status": "error",
// 				"error":  err.Error(),
// 				"output": string(output),
// 			}
// 		} else {
// 			results[pc.Name] = gin.H{
// 				"status": "success",
// 				"output": string(output),
// 			}
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{"results": results})
// }

package controllers

import (
	"fmt"
	"net/http"
	"os/exec"
	"wakeonlan/config"

	"github.com/gin-gonic/gin"
)

type PCs struct {
	Name         string
	IPAddress    string
	ComputerName string
	Username     string
	Password     string
}

func ExitSEBHandler(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT name, ip_address, computer_name, username, password
		FROM pcs where status = 0
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed: " + err.Error()})
		return
	}
	defer rows.Close()

	var pcs []PC
	for rows.Next() {
		var pc PC
		if err := rows.Scan(&pc.Name, &pc.IPAddress, &pc.ComputerName, &pc.Username, &pc.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB scan failed: " + err.Error()})
			return
		}
		pcs = append(pcs, pc)
	}

	results := make(map[string]interface{})

	for _, pc := range pcs {
		// Build the full PowerShell script
		// NOTE: We are using -Authentication Negotiate here instead of CredSSP
		fullScript := fmt.Sprintf(`
$secpasswd = ConvertTo-SecureString "%s" -AsPlainText -Force
$cred = New-Object System.Management.Automation.PSCredential ("%s", $secpasswd)
Invoke-Command -ComputerName "%s" -Credential $cred -Authentication Negotiate -ScriptBlock {
    Write-Host "Checking for SEB process..."
    $proc = Get-Process -Name "SafeExamBrowser" -ErrorAction SilentlyContinue
    if ($proc) {
        Write-Host "Stopping SEB..."
        Stop-Process -Id $proc.Id -Force
        Write-Host "Removing RecoveryState..."
        Remove-Item "C:\ProgramData\SafeExamBrowser\RecoveryState.seb" -Force -ErrorAction SilentlyContinue
        Write-Host "Shutting down..."
        shutdown /s /f /t 0
    } else {
        Write-Host "SEB process not found. No action taken."
    }
}
`, pc.Password, pc.Username, pc.IPAddress)

		// Execute the script
		cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", fullScript)
		output, err := cmd.CombinedOutput()

		if err != nil {
			results[pc.Name] = gin.H{
				"status": "error",
				"error":  err.Error(),
				"output": string(output),
			}
		} else {
			results[pc.Name] = gin.H{
				"status": "success",
				"output": string(output),
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
