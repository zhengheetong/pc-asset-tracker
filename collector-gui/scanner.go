package main

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/yusufpapurcu/wmi"
)

// PCSpecs holds the hardware data
type PCSpecs struct {
	OS         string `json:"os"`
	CPU        string `json:"cpu"`
	RAMTotal   string `json:"ramTotal"`
	RAMModules string `json:"ramModules"`
	Disks      string `json:"disks"`
	Serial     string `json:"serial"`
	Tag1       string `json:"tag1"`
	Tag2       string `json:"tag2"`
	Tag3       string `json:"tag3"`
}

// WMI Structs
type Win32_OperatingSystem struct {
	Caption string
}

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

	// OS
	var osData []Win32_OperatingSystem
	queryOS := wmi.CreateQuery(&osData, "")
	if err := wmi.Query(queryOS, &osData); err == nil && len(osData) > 0 {
		// This will grab things like "Microsoft Windows 11 Pro"
		specs.OS = osData[0].Caption
	}

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

	// Inside ScanHardware() in scanner.go
	var drives []Win32_DiskDrive
	// Remove the "WHERE" clause temporarily to see if your disks appear
	queryDisks := wmi.CreateQuery(&drives, "")

	if err := wmi.Query(queryDisks, &drives); err == nil && len(drives) > 0 {
		var validDisks []string
		for _, d := range drives {
			// Skip tiny partitions or virtual drives (less than 1GB)
			sizeGB := d.Size / (1000 * 1000 * 1000)

			// Also skip USB drives by checking the model or interface
			isUSB := strings.Contains(strings.ToLower(d.InterfaceType), "usb") ||
				strings.Contains(strings.ToLower(d.Model), "usb")

			if sizeGB > 0 && !isUSB {
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
