package checker

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Default server status URL
const DefaultURL string = "http://live.albiononline.com/status.txt"

// Default interval of status checking
const defaultInterval time.Duration = 1 * time.Minute

// Checker is a struct used for periodical server status checks.
type Checker struct {
	Client     *http.Client
	lastStatus StatusMessage
	Interval   time.Duration
	ticker     *time.Ticker
	c          chan bool
	closed     bool
	// Changes channel contains StatusMessages if there was server status change.
	Changes chan StatusMessage

	sync.RWMutex
}

// NewChecker returns default Checker
func NewChecker() *Checker {
	return &Checker{
		Client: &http.Client{
			Timeout: 3 * time.Second,
		},
		Changes:  make(chan StatusMessage, 1),
		Interval: defaultInterval,
	}
}

// GetStatus gets current status of albion server
func (c *Checker) GetStatus() (message StatusMessage, err error) {
	request, err := http.NewRequest(http.MethodGet, DefaultURL, nil)
	if err != nil {
		return
	}
	request.Header.Set("User-Agent", "Albion status checker")
	resp, err := c.Client.Do(request)
	if err != nil {
		if timerr, ok := err.(*url.Error); ok && timerr.Timeout() {
			message = StatusMessage{
				Status:  "timeout",
				Message: "Connection is timed out. Possibly service outage or DDoS.",
			}
			err = nil
			return
		}
		log.Debugf("Error occurred getting status.txt")
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	/* The server is sending a UTF-8 text string with a Byte Order Mark (BOM).
	The BOM identifies that the text is UTF-8 encoded, but it should be removed before decoding.
	https://stackoverflow.com/q/31398044
	Line breaks should be replaced with whitespaces because JSON standard doesn't allow this in strings. */
	content = bytes.ReplaceAll(bytes.TrimSpace((bytes.TrimPrefix(content, []byte("\xef\xbb\xbf")))), []byte{'\n'}, []byte{' '})

	err = json.Unmarshal(content, &message)
	if err == nil {
		log.Debugf("Retrieved server status: %s", message)
	}
	return
}

// CheckStatus checks current status of server and if it is different from last checking status sends change to Changes
func (c *Checker) CheckStatus() error {
	if c.Closed() {
		return errors.New("Checker is closed")
	}
	current, err := c.GetStatus()
	if err != nil {
		return err
	}
	c.Lock()
	if c.lastStatus.Status != current.Status {
		c.lastStatus = current
		c.Unlock()
		c.Changes <- current
		return nil
	}
	c.Unlock()
	return nil
}

func (c *Checker) loop() {
	for {
		select {
		case <-c.ticker.C:
			log.Debug("Loop")
			err := c.CheckStatus()
			if err != nil {
				log.Debugf("Error occurred while retrieving server status: %s", err)
			}
		case <-c.c:
			return
		}
	}
}

// Start starts checker goroutine
func (c *Checker) Start() error {
	if c.Closed() {
		return errors.New("Checker is closed")
	}
	c.c = make(chan bool)
	c.ticker = time.NewTicker(c.Interval)
	log.Debug("Starting looping goroutine")
	go c.loop()
	return nil
}

// Stop stops checker Ticker and closes Changes channel
func (c *Checker) Stop() error {
	if c.c == nil {
		return errors.New("Checker is not running")
	}
	if c.Closed() {
		return errors.New("Checker has been already stopped")
	}
	c.Lock()
	defer c.Unlock()
	c.ticker.Stop()
	c.c <- true
	close(c.Changes)
	return nil

}

// Closed returns current checker channel status
func (c *Checker) Closed() bool {
	c.RLock()
	defer c.RUnlock()
	return c.closed
}

// LastStatus returns latest StatusMessage
func (c *Checker) LastStatus() StatusMessage {
	c.RLock()
	defer c.RUnlock()
	return c.lastStatus
}
