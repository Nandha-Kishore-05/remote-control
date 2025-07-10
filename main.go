// package main

// import (
// 	"encoding/hex"
// 	"fmt"
// 	"net"
// 	"os/exec"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// var targetMACs = []string{
// 	"08:BF:B8:70:9E:23",
// 	"08:BF:B8:70:9F:3B",
// 	"08:BF:B8:70:9E:7A",
// }

// const broadcastIP = "255.255.255.255:9"

// func CreateMagicPacket(macAddr string) ([]byte, error) {
// 	macAddrClean := strings.ReplaceAll(macAddr, ":", "")
// 	macBytes, err := hex.DecodeString(macAddrClean)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid MAC address: %v", err)
// 	}

// 	packet := make([]byte, 6+16*6)
// 	for i := 0; i < 6; i++ {
// 		packet[i] = 0xFF
// 	}
// 	for i := 0; i < 16; i++ {
// 		copy(packet[6+i*6:], macBytes)
// 	}
// 	return packet, nil
// }

// func SendMagicPacket(mac string) error {
// 	packet, err := CreateMagicPacket(mac)
// 	if err != nil {
// 		return err
// 	}
// 	conn, err := net.Dial("udp", broadcastIP)
// 	if err != nil {
// 		return fmt.Errorf("error creating UDP connection: %v", err)
// 	}
// 	defer conn.Close()
// 	_, err = conn.Write(packet)
// 	if err != nil {
// 		return fmt.Errorf("error sending magic packet: %v", err)
// 	}
// 	return nil
// }

// func ExecutePowerShellScript() (string, error) {
// 	psScript := `
// $secpasswd = ConvertTo-SecureString "ITinfra@123" -AsPlainText -Force
// $cred = New-Object System.Management.Automation.PSCredential ("DESKTOP-B5VLFP5\cloud_70", $secpasswd)

// Invoke-Command -ComputerName DESKTOP-B5VLFP5 -Credential $cred -ScriptBlock {
//     $StartupPath = [Environment]::GetFolderPath("Startup")
//     $Shortcut = "$StartupPath\SafeExamBrowser.lnk"
//     $WshShell = New-Object -ComObject WScript.Shell
//     $ShortcutObject = $WshShell.CreateShortcut($Shortcut)
//     $ShortcutObject.TargetPath = "C:\Program Files (x86)\SafeExamBrowser\SEB.exe"
//     $ShortcutObject.Save()
// }
// `
// 	cmd := exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-Command", psScript)
// 	output, err := cmd.CombinedOutput()
// 	return string(output), err
// }

// func main() {
// 	router := gin.Default()

// 	router.GET("/poweron", func(c *gin.Context) {
// 		results := make(map[string]string)

// 		for _, mac := range targetMACs {
// 			err := SendMagicPacket(mac)
// 			if err != nil {
// 				results[mac] = fmt.Sprintf("Failed: %v", err)
// 			} else {
// 				results[mac] = "Magic packet sent"
// 			}
// 		}

// 		psOutput, psErr := ExecutePowerShellScript()
// 		if psErr != nil {
// 			results["PowerShell"] = fmt.Sprintf("Error: %v\nOutput:\n%s", psErr, psOutput)
// 		} else {
// 			results["PowerShell"] = fmt.Sprintf("Script executed successfully.\nOutput:\n%s", psOutput)
// 		}

// 		c.JSON(200, gin.H{
// 			"results": results,
// 		})
// 	})

// 	router.Run(":8080")
// }

// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"os/exec"

// 	"github.com/gin-gonic/gin"
// )

// type ShutdownRequest struct {
// 	TargetIPs []string `json:"target_ips" binding:"required"`
// }

// func main() {
// 	router := gin.Default()

// 	router.POST("/shutdown", func(c *gin.Context) {
// 		var req ShutdownRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// To collect per-IP results
// 		results := make(map[string]interface{})

// 		for _, ip := range req.TargetIPs {
// 			// Build shutdown command for each IP
// 			cmd := exec.Command(
// 				"shutdown",
// 				"/s",
// 				"/f",
// 				"/t", "0",
// 				"/m", fmt.Sprintf("\\\\%s", ip),
// 			)

// 			output, err := cmd.CombinedOutput()
// 			if err != nil {
// 				results[ip] = gin.H{
// 					"status": "error",
// 					"error":  err.Error(),
// 					"output": string(output),
// 				}
// 			} else {
// 				results[ip] = gin.H{
// 					"status": "success",
// 					"output": string(output),
// 				}
// 			}
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"results": results,
// 		})
// 	})

//		router.Run(":8080")
//	}

package main

import (
	"wakeonlan/config"
	turnoff "wakeonlan/controllers/TurnOffPc"
	turnon "wakeonlan/controllers/TurnOnPc"
	controllers "wakeonlan/controllers/safeExam"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	router.GET("/poweron", turnon.PowerOnHandler)
	router.POST("/shutdown", turnoff.ShutdownHandler)
	router.GET("/launch-seb", controllers.LaunchSEBHandler)
	router.Run(":8080")
}
