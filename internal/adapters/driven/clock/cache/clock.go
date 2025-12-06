package cache

import (
	"os"
	"time"
)

type CacheClock struct {
	filePath string
}

func NewClock(path string) CacheClock {
	return CacheClock{filePath: path}
}

func (c CacheClock) Current() time.Time {
	currentTime, err := c.parseFromFile()

	if err != nil {
		panic(err)
	}

	return currentTime
}

func (c CacheClock) Set(t time.Time) error {
	return os.WriteFile(c.filePath, []byte(t.Format(time.DateTime)), 0666)
}

func (c CacheClock) Path() string {
	return c.filePath
}

func (c CacheClock) parseFromFile() (time.Time, error) {
	err := c.createIfNotExist()

	if err != nil {
		return time.Time{}, err
	}

	data, err := os.ReadFile(c.filePath)

	if err != nil {
		return time.Time{}, err
	}

	timeString := string(data)

	if timeString == "" {
		return time.Now(), nil
	}
	parsed, err := time.Parse(time.DateTime, timeString)

	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

func (c CacheClock) createIfNotExist() error {
	_, err := os.Stat(c.filePath)
	if err != nil{
		if os.IsNotExist(err) {
			file, err := os.Create(c.filePath)

			if err != nil {
				return err
			}
			defer file.Close()
			return nil
		}
		return err
	}

	return nil;
}
