package checker

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Default server status url
const url string = "http://live.albiononline.com/status.txt"

// Default interval of status checking
const defaultInterval time.Duration = 3 * time.Minute

type Checker struct {
	Client     *http.Client
	lastStatus StatusMessage
	Ticker     *time.Ticker
	closed     bool
	// Changes channel contains StatusMessages if there was server status change.
	Changes chan StatusMessage

	sync.RWMutex
}

// NewChecker returns default Checker
func NewChecker() *Checker {
	return &Checker{
		Client:  http.DefaultClient,
		Changes: make(chan StatusMessage),
		Ticker:  time.NewTicker(defaultInterval),
		closed:  false,
	}
}

// GetStatus gets current status of albion server
func (c *Checker) GetStatus() (message StatusMessage, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	request.Header.Set("User-Agent", "Albion status checker")
	resp, err := c.Client.Do(request)
	if err != nil {
		log.Debugf("Error occured getting status.txt")
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	/* The server is sending a UTF-8 text string with a Byte Order Mark (BOM).
	The BOM identifies that the text is UTF-8 encoded, but it should be removed before decoding.
	https://stackoverflow.com/q/31398044 */
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
	defer c.Unlock()
	if c.lastStatus.Status != current.Status {
		c.lastStatus = current
		c.Changes <- current
	}
	return nil
}
func (c *Checker) loop() {
	for range c.Ticker.C {
		log.Debug("Loop")
		err := c.CheckStatus()
		if err != nil {
			log.Debugf("Error occured while retrieving server status: %s", err)
		}
	}
}

// Start starts checker goroutine
func (c *Checker) Start() error {
	if c.Closed() {
		return errors.New("Checker is closed")
	}
	log.Debug("Starting looping goroutine")
	go c.loop()
	return nil
}

// Stop stops checker Ticker and closes Changes channel
func (c *Checker) Stop() error {
	if !c.Closed() {
		c.Lock()
		defer c.Unlock()
		c.Ticker.Stop()
		close(c.Changes)
		c.closed = true
		return nil
	}
	return errors.New("Checker has been already stopped")
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
