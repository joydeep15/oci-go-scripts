package main

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common/auth"
	"github.com/oracle/oci-go-sdk/loggingingestion"
	"time"

	"flag"
	"github.com/google/uuid"
	"github.com/oracle/oci-go-sdk/common"
)

func main() {
	// Create a context with a default timeout
	logOcid := flag.String("l", "xyz", "log id where we want to write (ocid1.log.iad.xyz)")
	message := flag.String("m", "", "optional message prefix")
	flag.Parse()
	fmt.Printf("Writing to logocid: %s \n", *logOcid)
	fmt.Printf("messageprefix: %s \n", *message)
	ctx := context.Background()

	// Create a client for the Logging Management service
	instancePrincipalProvider, err := auth.InstancePrincipalConfigurationProvider()
	if err != nil {
		fmt.Println("error while creating instance principal provider %w", err)
	}
	client, err := loggingingestion.NewLoggingClientWithConfigurationProvider(instancePrincipalProvider)
	if err != nil {
		fmt.Println("Failed to create Logging Management client:", err)
		return
	}

	// Define the log entry to be written
	fmt.Printf("writing message: %s\n", *message)
	logEntry := []loggingingestion.LogEntry{{
		Data: message,
		Id:   common.String(uuid.New().String()),
		Time: nil,
	}}

	// Write the log entry to the specified OCI log ID
	resp, err := client.PutLogs(ctx,
		loggingingestion.PutLogsRequest{
			PutLogsDetails: loggingingestion.PutLogsDetails{
				LogEntryBatches: []loggingingestion.LogEntryBatch{{Defaultlogentrytime: &common.SDKTime{Time: time.Now()},
					Entries: logEntry,
					Source:  common.String("my go client"),
					Type:    common.String("my custom log")}},
				Specversion: common.String("1.0")},
			TimestampOpcAgentProcessing: &common.SDKTime{Time: time.Now()},
			LogId:                       logOcid},
	)

	if err != nil {
		fmt.Println("Failed to write log entry:", err)
		return
	}
	fmt.Printf("opcreqid: %s \n", *resp.OpcRequestId)
	fmt.Println("Log entry written successfully")
}
