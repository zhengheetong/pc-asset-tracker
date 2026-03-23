package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// UploadToGoogleSheets pushes the PCSpecs to your cloud spreadsheet
func UploadToGoogleSheets(specs PCSpecs) error {
	ctx := context.Background()

	// 1. Path to your service account key file
	credentialsFile := "service-account.json"

	// 2. Your Spreadsheet ID (Copy this from the URL of your sheet)
	spreadsheetId := "1WK_dJKXyquSS0xivuej8Ic8WxPDBLpYsSvUVVFP3eM8"

	// 3. Initialize the Sheets service
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	// 4. Prepare the row data
	// Order: Timestamp | Serial | CPU | Total RAM | Modules | Disks
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	row := []interface{}{
		timestamp,
		specs.Serial,
		specs.OS,
		specs.CPU,
		specs.RAMTotal,
		specs.RAMModules,
		specs.Disks,
		specs.Tag1,
		specs.Tag2,
		specs.Tag3,
	}

	rb := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	// 5. Append the row to the bottom of "Sheet1"
	writeRange := "Sheet1!A1"
	_, err = srv.Spreadsheets.Values.Append(spreadsheetId, writeRange, rb).
		ValueInputOption("RAW").
		InsertDataOption("INSERT_ROWS").
		Context(ctx).
		Do()

	return err
}
