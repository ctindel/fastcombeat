package beater

import (
	"fmt"
	"time"
	"net/http"
	"io"
	"os"
	"strconv"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/ctindel/fastcombeat/config"
	"github.com/ctindel/fastcombeat/fast"
	"github.com/ctindel/fastcombeat/format"
	"github.com/ctindel/fastcombeat/meters"
)

const debugK = "fastcombeat"

type Fastcombeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

func calculateBandwidth(urls []string, bytesPerSec *uint64) (err error) {
	client := &http.Client{}
	count := uint64(len(urls))

	primaryBandwidthReader := meters.BandwidthMeter{}
	bandwidthMeter := meters.BandwidthMeter{}
	ch := make(chan *copyResults, 1)
	bytesToRead := uint64(0)
	completed := uint64(0)

	for i := uint64(0); i < count; i++ {
		// Create the HTTP request
		request, err := http.NewRequest("GET", urls[i], nil)
		if err != nil {
			return err
		}
		request.Header.Set("User-Agent", "fastcombeat 1.0")

		// Get the HTTP Response
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		// Set information for the leading index
		if i == 0 {
			// Try to get content length
			contentLength := response.Header.Get("Content-Length")
			calculatedLength, err := strconv.Atoi(contentLength)
			if err != nil {
				calculatedLength = 26214400
			}
			bytesToRead = uint64(calculatedLength)
			logp.Debug(debugK, "Download Size=%d", bytesToRead)

			tapMeter := io.TeeReader(response.Body, &primaryBandwidthReader)
			go asyncCopy(i, ch, &bandwidthMeter, tapMeter)
		} else {
			// Start reading
			go asyncCopy(i, ch, &bandwidthMeter, response.Body)
		}

	}

	logp.Debug(debugK, "Estimating current download speed")
	for {
		select {
		case results := <-ch:
			if results.err != nil {
				logp.Err("%s", results.err)
				os.Exit(1)
			}

			completed++
			logp.Debug(debugK, "%s - %s",
				format.BitsPerSec(bandwidthMeter.Bandwidth()),
				format.Percent(primaryBandwidthReader.BytesRead(), bytesToRead))
			logp.Debug(debugK, "Completed in %.1f seconds", bandwidthMeter.Duration().Seconds())
			*bytesPerSec = uint64(primaryBandwidthReader.Bandwidth())
			return nil
		case <-time.After(100 * time.Millisecond):
			logp.Info("%s - %s",
				format.BitsPerSec(bandwidthMeter.Bandwidth()),
				format.Percent(primaryBandwidthReader.BytesRead(), bytesToRead))
		}
	}
}

type copyResults struct {
	index        uint64
	bytesWritten uint64
	err          error
}

func asyncCopy(index uint64, channel chan *copyResults, writer io.Writer, reader io.Reader) {
	bytesWritten, err := io.Copy(writer, reader)
	channel <- &copyResults{index, uint64(bytesWritten), err}
}

func sumArr(array []uint64) (sum uint64) {
	for i := 0; i < len(array); i++ {
		sum = sum + array[i]
	}
	return
}
// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Fastcombeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Fastcombeat) Run(b *beat.Beat) error {
	logp.Info("fastcombeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}


		count := uint64(3)
		bytesPerSec := uint64(0)

		fast.UseHTTPS = true
		urls := fast.GetDlUrls(count)
		logp.Info("Got %d from fast service", len(urls))
		if len(urls) == 0 {
			logp.Info("Using fallback endpoint")
			urls = append(urls, fast.GetDefaultURL())
		}

		err := calculateBandwidth(urls, &bytesPerSec)
		if err != nil {
			logp.Err("%s", err)
		}

		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"bytespersecond": bytesPerSec,
			},
		}
		bt.client.Publish(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Fastcombeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
