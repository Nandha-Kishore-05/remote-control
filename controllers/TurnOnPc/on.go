// package controllers

// import (
// 	"encoding/hex"
// 	"fmt"
// 	"net"
// 	"os/exec"
// 	"strings"
// 	"wakeonlan/config"

// 	"github.com/gin-gonic/gin"
// )

// type PC struct {
// 	ID           int
// 	Name         string
// 	MACAddress   string
// 	IPAddress    string
// 	ComputerName string
// 	Username     string
// 	Password     string
// 	Status       int

// }

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
// 	conn, err := net.Dial("udp", "255.255.255.255:9")
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

// func PowerOnHandler(c *gin.Context) {
// 	// Fetch all PCs
// 	rows, err := config.DB.Query(`
// 		SELECT mac_address, computer_name, username, password
// 		FROM pcs
// 	`)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": "DB query failed: " + err.Error()})
// 		return
// 	}
// 	defer rows.Close()

// 	var pcs []PC
// 	for rows.Next() {
// 		var pc PC
// 		if err := rows.Scan(&pc.MACAddress, &pc.ComputerName, &pc.Username, &pc.Password); err != nil {
// 			c.JSON(500, gin.H{"error": "DB scan failed: " + err.Error()})
// 			return
// 		}
// 		pcs = append(pcs, pc)
// 	}

// 	results := make(map[string]string)

// 	for _, pc := range pcs {

// 		if err := SendMagicPacket(pc.MACAddress); err != nil {
// 			results[pc.MACAddress] = fmt.Sprintf("WOL failed: %v", err)
// 		} else {
// 			results[pc.MACAddress] = "Magic packet sent"
// 		}

// 		psScript := fmt.Sprintf(`
// $secpasswd = ConvertTo-SecureString "%s" -AsPlainText -Force
// $cred = New-Object System.Management.Automation.PSCredential ("%s\%s", $secpasswd)
// Invoke-Command -ComputerName %s -Credential $cred -ScriptBlock {
//     $StartupPath = [Environment]::GetFolderPath("Startup")
//     $Shortcut = "$StartupPath\SafeExamBrowser.lnk"
//     $WshShell = New-Object -ComObject WScript.Shell
//     $ShortcutObject = $WshShell.CreateShortcut($Shortcut)
//     $ShortcutObject.TargetPath = "C:\Program Files (x86)\SafeExamBrowser\SEB.exe"
//     $ShortcutObject.Save()
// }
// `, pc.Password, pc.ComputerName, pc.Username, pc.ComputerName)

// 		cmd := exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-Command", psScript)
// 		output, err := cmd.CombinedOutput()
// 		if err != nil {
// 			results[pc.ComputerName] = fmt.Sprintf("PowerShell error: %v\nOutput:\n%s", err, output)
// 		} else {
// 			results[pc.ComputerName] = fmt.Sprintf("PowerShell script executed successfully.\nOutput:\n%s", output)
// 		}
// 	}

//		c.JSON(200, gin.H{"results": results})
//	}

// package controllers

// import (
// 	"encoding/hex"
// 	"fmt"
// 	"net"
// 	"os/exec"

// 	//"os/exec"
// 	"strings"
// 	"wakeonlan/config"

// 	"github.com/gin-gonic/gin"

// )

// type PC struct {
// 	ID           int
// 	Name         string
// 	MACAddress   string
// 	IPAddress    string
// 	ComputerName string
// 	Username     string
// 	Password     string
// 	Status       int
// }

// // CreateMagicPacket prepares a WOL packet
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

// // SendMagicPacket sends the packet over UDP broadcast
// func SendMagicPacket(mac string) error {
// 	packet, err := CreateMagicPacket(mac)
// 	if err != nil {
// 		return err
// 	}
// 	conn, err := net.Dial("udp", "255.255.255.255:9")
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

// // PowerOnHandler wakes up PCs and schedules SafeExamBrowser at startup
// func PowerOnHandler(c *gin.Context) {
//   rows, err := config.DB.Query(`
//     SELECT mac_address, ip_address, computer_name, username, password
//     FROM pcs
//   `)
//   if err != nil {
//     c.JSON(500, gin.H{"error": "DB query failed: " + err.Error()})
//     return
//   }
//   defer rows.Close()

//   var pcs []PC
//   for rows.Next() {
//     var pc PC
//     if err := rows.Scan(&pc.MACAddress, &pc.IPAddress, &pc.ComputerName, &pc.Username, &pc.Password); err != nil {
//       c.JSON(500, gin.H{"error": "DB scan failed: " + err.Error()})
//       return
//     }
//     pcs = append(pcs, pc)
//   }

//   results := make(map[string]string)

//   for _, pc := range pcs {
//     // Send WOL
//     if err := SendMagicPacket(pc.MACAddress); err != nil {
//       results[pc.MACAddress] = fmt.Sprintf("WOL failed: %v", err)
//     } else {
//       results[pc.MACAddress] = "Magic packet sent"
//     }

//     // PowerShell script
//   psScript := fmt.Sprintf(
//   `$remotePC = "%s"; $remoteUser = "%s\%s"; $remotePass = "%s"; `+
//     `$securePass = ConvertTo-SecureString $remotePass -AsPlainText -Force; `+
//     `$cred = New-Object System.Management.Automation.PSCredential ($remoteUser, $securePass); `+
//     `Invoke-Command -ComputerName $remotePC -Credential $cred -ScriptBlock { `+
//     `Get-Process }`,
//   pc.IPAddress, pc.ComputerName, pc.Username, pc.Password,
// )

//     cmd := exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-Command", psScript)
//     output, err := cmd.CombinedOutput()
//     if err != nil {
//       results[pc.ComputerName] = fmt.Sprintf("PowerShell error: %v\nOutput:\n%s", err, string(output))
//     } else {
//       results[pc.ComputerName] = fmt.Sprintf("Scheduled task created successfully.\nOutput:\n%s", string(output))
//     }
//   }

//   c.JSON(200, gin.H{"results": results})
// }

package controllers

import (

	"encoding/hex"
	"fmt"
	"net"


	"strings"

	"wakeonlan/config"

	"github.com/gin-gonic/gin"
)

type PC struct {
	ID           int
	Name         string
	MACAddress   string
	IPAddress    string
	ComputerName string
	Username     string
	Password     string
	Status       int
}

// CreateMagicPacket prepares a WOL packet
func CreateMagicPacket(macAddr string) ([]byte, error) {
	macAddrClean := strings.ReplaceAll(macAddr, ":", "")
	macBytes, err := hex.DecodeString(macAddrClean)
	if err != nil {
		return nil, fmt.Errorf("invalid MAC address: %v", err)
	}

	packet := make([]byte, 6+16*6)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 0; i < 16; i++ {
		copy(packet[6+i*6:], macBytes)
	}
	return packet, nil
}

// SendMagicPacket sends the packet over UDP broadcast
func SendMagicPacket(mac string) error {
	packet, err := CreateMagicPacket(mac)
	if err != nil {
		return err
	}
	conn, err := net.Dial("udp", "255.255.255.255:9")
	if err != nil {
		return fmt.Errorf("error creating UDP connection: %v", err)
	}
	defer conn.Close()
	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("error sending magic packet: %v", err)
	}
	return nil
}

func PowerOnHandler(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT mac_address, ip_address, computer_name, username, password
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
		if err := rows.Scan(&pc.MACAddress, &pc.IPAddress, &pc.ComputerName, &pc.Username, &pc.Password); err != nil {
			c.JSON(500, gin.H{"error": "DB scan failed: " + err.Error()})
			return
		}
		pcs = append(pcs, pc)
	}

	results := make(map[string]string)

	for _, pc := range pcs {
		// 1. Wake-on-LAN
		if err := SendMagicPacket(pc.MACAddress); err != nil {
			results[pc.IPAddress] = fmt.Sprintf("WOL failed: %v", err)
			continue
		}

		// 2. PowerShell command to run remotely

	}

	c.JSON(200, gin.H{"results": results})
}