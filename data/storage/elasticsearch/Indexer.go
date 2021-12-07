package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/tkodyl/vineguard/configuration"
	"github.com/tkodyl/vineguard/data/collection/pm"
	"log"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

type Indexer struct {
	Config          configuration.Config
	countSuccessful uint64
}

func (indexer *Indexer) Do(records []*pm.WeatherRecord) {
	log.Printf("About to start of indexing %d\n", len(records))

	bi, err := indexer.createBulkIndexer()
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	start := time.Now().UTC()
	for _, record := range records {
		data, err := json.Marshal(record)
		if err != nil {
			log.Fatalf("Cannot encode article from %s %s: %s", record.Date, record.Time, err)
		}

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "create", // index
				Body:   bytes.NewReader(data),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&indexer.countSuccessful, 1)
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			})
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
	dur := time.Since(start)
	indexer.printReport(bi, dur)
	log.Printf("Total number of sent messages: %d\n", indexer.countSuccessful)
}

func (indexer *Indexer) createBulkIndexer() (esutil.BulkIndexer, error) {
	retryBackoff := backoff.NewExponentialBackOff()

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:     []string{indexer.Config.Elasticsearch.Url},
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexer.Config.Elasticsearch.Index, // The default index name
		Client:        es,                                 // The Elasticsearch client
		NumWorkers:    runtime.NumCPU(),                   // The number of worker goroutines
		FlushBytes:    int(5e6),                           // The flush threshold in bytes
		FlushInterval: 30 * time.Second,                   // The periodic flush interval
	})
	return bi, err
}

func (indexer *Indexer) printReport(bi esutil.BulkIndexer, dur time.Duration) {
	biStats := bi.Stats()

	log.Println(strings.Repeat("▔", 65))

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%d] documents with [%d] errors in %d (%d docs/sec)",
			int64(biStats.NumFlushed),
			int64(biStats.NumFailed),
			dur.Truncate(time.Millisecond),
			int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed)),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%d] documents in %d (%d docs/sec)",
			int64(biStats.NumFlushed),
			dur.Truncate(time.Millisecond),
			int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed)),
		)
	}
}
