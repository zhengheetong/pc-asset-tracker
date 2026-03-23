package main

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/yusufpapurcu/wmi"
)

// PCSpecs holds the hardware data
type PCSpecs struct {
	CPU        string
	RAMTotal   string
	RAMModules string
	Disks      string
	Serial     string
	// NEW: Meta Tags for categorization
	Tag1 string
	Tag2 string
	Tag3 string
}

// WMI Structs
type Win32_PhysicalMemory struct {
	Manufacturer, PartNumber, DeviceLocator string
	Capacity                                uint64
	Speed                                   uint32
}

type Win32_DiskDrive struct {
	Model         string
	Size          uint64
	InterfaceType string
}

type Win32_BIOS struct{ SerialNumber string }
type Win32_ComputerSystemProduct struct{ IdentifyingNumber string }
type Win32_BaseBoard struct{ SerialNumber string }

func ScanHardware() (PCSpecs, error) {
	var specs PCSpecs

	// 1. CPU
	cpuInfo, _ := cpu.Info()
	if len(cpuInfo) > 0 {
		specs.CPU = cpuInfo[0].ModelName
	}

	// 2. Get RAM Info (Total + Individual Modules via WMI)
	// 2. Get RAM Info (Total + Individual Modules via WMI)
	var ramModules []Win32_PhysicalMemory
	qRam := wmi.CreateQuery(&ramModules, "")

	if err := wmi.Query(qRam, &ramModules); err == nil && len(ramModules) > 0 {
		var modules []string
		var physicalTotalBytes uint64

		for _, mod := range ramModules {
			physicalTotalBytes += mod.Capacity
			sizeGB := mod.Capacity / (1024 * 1024 * 1024)

			// --- NEW: Trim DeviceLocator ---
			// Examples: "Controller0-ChannelA" -> "A", "DIMM 1" -> "1"
			slot := mod.DeviceLocator
			slot = strings.ReplaceAll(slot, "Controller0-", "")
			slot = strings.ReplaceAll(slot, "Channel", "")
			slot = strings.ReplaceAll(slot, "DIMM", "")
			slot = strings.TrimSpace(strings.ReplaceAll(slot, "-", ""))

			// Add "Slot" prefix for better readability
			displaySlot := fmt.Sprintf("Slot %s", slot)
			// -------------------------------

			vendor := strings.TrimSpace(mod.Manufacturer)
			if vendor == "" || strings.ToLower(vendor) == "unknown" {
				vendor = "Generic"
			}

			part := strings.TrimSpace(mod.PartNumber)
			partStr := ""
			if part != "" && strings.ToLower(part) != "unknown" {
				partStr = fmt.Sprintf(" [%s]", part)
			}

			moduleStr := fmt.Sprintf("%s: %s %dGB %dMHz%s",
				displaySlot, // Use the trimmed slot name
				vendor,
				sizeGB,
				mod.Speed,
				partStr,
			)
			modules = append(modules, moduleStr)
		}

		specs.RAMTotal = fmt.Sprintf("%d GB", physicalTotalBytes/(1024*1024*1024))
		specs.RAMModules = strings.Join(modules, " | ")
	}

	// 3. FAST Disk Scan (Using WMI)
	var drives []Win32_DiskDrive
	// We filter out USB drives directly in the query to save time
	queryDisks := wmi.CreateQuery(&drives, "WHERE InterfaceType != 'USB'")

	if err := wmi.Query(queryDisks, &drives); err == nil && len(drives) > 0 {
		var validDisks []string
		for _, d := range drives {
			// Marketing Math: Manufacturers use 1000^3
			// 1,000,000,000,000 / (1000^3) = 1000 GB (1TB)
			sizeGB := d.Size / (1000 * 1000 * 1000)

			if sizeGB > 0 {
				driveStr := fmt.Sprintf("%s (%dGB)", strings.TrimSpace(d.Model), sizeGB)
				validDisks = append(validDisks, driveStr)
			}
		}
		specs.Disks = strings.Join(validDisks, " | ")
	} else {
		specs.Disks = "Unknown Disks"
	}

	// 4. Serial Number
	specs.Serial = getSerialNumber()

	return specs, nil
}

func getSerialNumber() string {
	var bios []Win32_BIOS
	if err := wmi.Query(wmi.CreateQuery(&bios, ""), &bios); err == nil && len(bios) > 0 {
		if isValidSerial(bios[0].SerialNumber) {
			return strings.TrimSpace(bios[0].SerialNumber)
		}
	}

	var csp []Win32_ComputerSystemProduct
	if err := wmi.Query(wmi.CreateQuery(&csp, ""), &csp); err == nil && len(csp) > 0 {
		if isValidSerial(csp[0].IdentifyingNumber) {
			return strings.TrimSpace(csp[0].IdentifyingNumber)
		}
	}

	var board []Win32_BaseBoard
	if err := wmi.Query(wmi.CreateQuery(&board, ""), &board); err == nil && len(board) > 0 {
		if isValidSerial(board[0].SerialNumber) {
			return strings.TrimSpace(board[0].SerialNumber)
		}
	}

	return "Unknown Serial"
}

func isValidSerial(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	bad := []string{"", "none", "unknown", "default string", "to be filled by o.e.m.", "system serial number", "00000000"}
	for _, b := range bad {
		if s == b {
			return false
		}
	}
	return true
}
