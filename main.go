package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/cheggaaa/pb"
	"google.golang.org/api/iterator"
)

var header string = "device_id,destination_ip,mo_data_in_mb,mt_data_in_mb\n"

func main() {
	projectID := "ttm-aersight-b4c15e1f"

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	rows, err := query(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	if err := writeResultsToCSV(os.Stdout, rows); err != nil {
		log.Fatal(err)
	}
}

// query returns a row iterator suitable for reading query results.
func query(ctx context.Context, client *bigquery.Client) (*bigquery.RowIterator, error) {
	fmt.Println("Requesting Data from BigQuery...")

	query := client.Query(
		`select 
		device_id,  
		destination_ip, 
		cast(sum(uplink_data/1024/1024) as numeric) as mo_data_in_mb, 
    	cast(sum(downlink_data/1024/1024) as numeric) as mt_data_in_mb,  
	from 
	` + "`ttm-aersight-b4c15e1f.aersight_acp_ds.ipstream_flows`" + `
	WHERE 
		eventdate >= TIMESTAMP(DATE_SUB(CURRENT_DATE(), interval 1 day))
    	and eventdate < timestamp(CURRENT_DATE())
	GROUP BY 
		device_id, destination_ip;`)
	return query.Read(ctx)
}

func writeResultsToCSV(w io.Writer, iter *bigquery.RowIterator) error {

	var nRows int
	var bar *pb.ProgressBar

	file, err := os.Create("result.csv")
	checkFatalError("Cannot create file", err)
	defer file.Close()

	file.WriteString(header)

	for {
		var row []bigquery.Value
		err := iter.Next(&row)
		if nRows == 0 {
			nRows = int(iter.TotalRows)
			fmt.Println("Number of Rows to process:", nRows)
			bar = pb.StartNew(nRows)
			defer bar.Finish()
		}

		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error iterating through results: %v", err)
		}
		szTemp := ""
		for i := range row {
			if i > 0 {
				szTemp += ","
			}
			szTemp += fmt.Sprint(row[i])

		}
		szTemp += "\n"
		_, err = file.WriteString(szTemp)
		checkFatalError("Cannot write to file", err)
		bar.Increment()
	}
}

func checkFatalError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
