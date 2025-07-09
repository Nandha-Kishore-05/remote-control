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
	// Fetch all PCs from DB
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
	
		cmd := exec.Command(
			"shutdown",
			"/s",
			"/f",
			"/t", "0",
			"/m", fmt.Sprintf("\\\\%s", pc.IPAddress),
		)

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

	c.JSON(200, gin.H{"results": results})
}
