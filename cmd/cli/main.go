package main

import (
	"context"
	"encoding/csv"
	"flag"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shumon84/fetch-twitter-icon/pkg/icon"
)

func main() {
	var (
		outputPath   = flag.String("o", "output/", "output path")
		token        = flag.String("t", "", "Twitter API bearer token")
		inputCSVPath = flag.String("i", "input.csv", "input csv file")
	)
	flag.Parse()
	inputFile, err := os.Open(*inputCSVPath)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	inputCSV := csv.NewReader(inputFile)
	if _, err := inputCSV.Read(); err != nil { // CSVの 1 行目はスキップする
		log.Fatal(err)
	}
	client := icon.NewFetchClient(*token, http.DefaultClient)
	for {
		record, err := inputCSV.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(record) == 0 {
			log.Fatal("invalid csv")
		}
		twitterID := record[0]
		download(client, twitterID, *outputPath)
	}
}

func download(client icon.FetchClient, twitterID string, outputPath string) {
	iconImage, err := client.Fetch(context.Background(), twitterID)
	if err != nil {
		log.Fatalln(err)
	}
	outputFile, err := os.Create(filepath.Join(outputPath, twitterID+".png"))
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	if err := png.Encode(outputFile, iconImage); err != nil {
		log.Fatal(err)
	}
}
