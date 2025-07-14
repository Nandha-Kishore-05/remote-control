package controllers

import (
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
func ShutdownHandler(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT name, ip_address, computer_name, username, password
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
		if err := rows.Scan(&pc.Name, &pc.IPAddress, &pc.ComputerName, &pc.Username, &pc.Password); err != nil {
			c.JSON(500, gin.H{"error": "DB scan failed: " + err.Error()})
			return
		}
		pcs = append(pcs, pc)
	}

	results := make(map[string]interface{})

	for _, pc := range pcs {
		// Step 1: Authenticate with net use
		netUseCmd := exec.Command("net", "use", fmt.Sprintf(`\\%s`, pc.IPAddress), fmt.Sprintf("/user:%s", pc.Username), pc.Password)
		netUseOutput, netUseErr := netUseCmd.CombinedOutput()

		if netUseErr != nil {
			results[pc.Name] = gin.H{
				"status": "net use failed",
				"error":  netUseErr.Error(),
				"output": string(netUseOutput),
			}
			continue // Skip shutdown if net use fails
		}

		// Step 2: Run shutdown command
		shutdownCmd := exec.Command(
			"shutdown",
			"/s",
			"/f",
			"/t", "0",
			"/m", fmt.Sprintf(`\\%s`, pc.IPAddress),
		)
		shutdownOutput, shutdownErr := shutdownCmd.CombinedOutput()

		if shutdownErr != nil {
			results[pc.Name] = gin.H{
				"status":  "shutdown failed",
				"error":   shutdownErr.Error(),
				"output":  string(shutdownOutput),
				"net_use": string(netUseOutput),
			}
		} else {
			results[pc.Name] = gin.H{
				"status":  "shutdown success",
				"output":  string(shutdownOutput),
				"net_use": string(netUseOutput),
			}
		}
	}

	c.JSON(200, gin.H{"results": results})
}
