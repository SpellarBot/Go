// Conf Loader
// @Author: Golion
// @Date: 2017.6

package utils

import (
	"os"
	"bufio"
	"io"
	"strings"
	"sync"
	"time"
	"fmt"
)

type ConfLoader struct {
	FileName       string
	OutdateSeconds int

	kvMap          map[string]string
	mutex          sync.RWMutex
	ticker         *time.Ticker
}

func NewConfLoader(fileName string, outdateSeconds int) (*ConfLoader, error) {
	confLoader := ConfLoader{
		FileName:       fileName,
		OutdateSeconds: outdateSeconds,
	}
	if err := confLoader.Init(); err != nil {
		return nil, err
	} else {
		return &confLoader, nil
	}
}

func (c *ConfLoader) Init() error {
	if len(c.FileName) > 0 {
		err := c.loadFile()
		if err != nil {
			return err
		} else {
			go c.loopUpdateCheck()
			return nil
		}
	} else {
		return fmt.Errorf("[ConfLoader][Init] Error! File Name Is Empty!")
	}
}

func (c *ConfLoader) loadFile() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !IsFileExist(c.FileName) {
		return fmt.Errorf("File Not Exist! fileName=[%v]", c.FileName)
	}
	f, err := os.Open(c.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	c.kvMap = make(map[string]string)
	br := bufio.NewReader(f)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.Replace(line, "\n", "", -1)
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		equal := -1
		for k, c := range line {
			if c == '=' {
				equal = k
				break
			}
		}
		if equal > 0 {
			key   := strings.TrimSpace(line[:equal])
			value := strings.TrimSpace(line[equal+1:])
			if len(key) > 0 && len(value) > 0 {
				c.kvMap[key] = value
			}
		}
	}
	return nil
}

func (c *ConfLoader) loopUpdateCheck() {
	if c.OutdateSeconds > 0 {
		c.ticker = time.NewTicker(time.Duration(c.OutdateSeconds) * time.Second)
		for _ = range c.ticker.C {
			c.loadFile()
		}
	}
}

func (c *ConfLoader) String(key string) string {
	if c.kvMap == nil {
		return ""
	}
	c.mutex.RLock()
	c.mutex.RUnlock()
	if val, ok := c.kvMap[key]; ok {
		return val
	}
	return ""
}

func (c *ConfLoader) Bool(key string) bool {
	if c.kvMap == nil {
		return false
	}
	c.mutex.RLock()
	c.mutex.RUnlock()
	if val, ok := c.kvMap[key]; ok {
		return (strings.ToLower(val) == "true") || (val == "1")
	}
	return false
}

func (c *ConfLoader) Int(key string) int {
	if c.kvMap == nil {
		return 0
	}
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if val, ok := c.kvMap[key]; ok {
		return Atoi(val)
	}
	return 0
}

func (c *ConfLoader) Int32(key string) int32 {
	if c.kvMap == nil {
		return int32(0)
	}
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if val, ok := c.kvMap[key]; ok {
		return Atoi32(val)
	}
	return int32(0)
}

func (c *ConfLoader) Int64(key string) int64 {
	if c.kvMap == nil {
		return int64(0)
	}
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if val, ok := c.kvMap[key]; ok {
		return Atoi64(val)
	}
	return int64(0)
}

func (c *ConfLoader) Float32(key string) float32 {
	if c.kvMap == nil {
		return float32(0)
	}
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if val, ok := c.kvMap[key]; ok {
		return Atof32(val)
	}
	return float32(0)
}

func (c *ConfLoader) Float64(key string) float64 {
	if c.kvMap == nil {
		return float64(0)
	}
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if val, ok := c.kvMap[key]; ok {
		return Atof64(val)
	}
	return float64(0)
}