package main

import (
	"encoding/csv"
	"os"

	"github.com/gocarina/gocsv"
)

func writeActions(out string, actions []Action) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = '\t'
	return gocsv.MarshalCSV(actions, w)
}

func writeOmitted(out string, omitted []Omit) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = '\t'
	return gocsv.MarshalCSV(omitted, w)
}

func writeRecordStream(out string, records <-chan Record) (done <-chan error, err error) {
	f, err := os.Create(out)
	if err != nil {
		return nil, err
	}
	completion := make(chan error)
	proxy := make(chan interface{})
	go func() {
		defer close(completion)
		defer f.Close()
		w := csv.NewWriter(f)
		w.Comma = '\t'
		completion <- gocsv.MarshalChan(proxy, w)
	}()
	go func() {
		defer close(proxy)
		for record := range records {
			proxy <- record
		}
	}()
	return completion, nil
}
