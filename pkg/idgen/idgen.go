package idgen

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch         = int64(1577836800000) // 2020-01-01 00:00:00 UTC
	machineIDBits = uint(10)
	sequenceBits  = uint(12)

	maxMachineID = int64(-1) ^ (int64(-1) << machineIDBits)
	maxSequence  = int64(-1) ^ (int64(-1) << sequenceBits)
)

type Snowflake struct {
	machineID int64
	mu        sync.Mutex
	lastTime  int64
	sequence  int64
}

var (
	instance *Snowflake
	once     sync.Once
)

func New(machineID int64) *Snowflake {
	if machineID < 0 || machineID > maxMachineID {
		panic("invalid machine ID")
	}

	return &Snowflake{
		machineID: machineID,
	}
}

func GetInstance() *Snowflake {
	once.Do(func() {
		instance = New(1) // 默认使用机器ID 1
	})
	return instance
}

func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano() / 1000000
	if now < s.lastTime {
		panic("clock moved backwards")
	}

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.lastTime {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now

	return (now-epoch)<<(machineIDBits+sequenceBits) |
		(s.machineID << sequenceBits) |
		s.sequence
}

func GenerateString() string {
	return fmt.Sprintf("%d", GetInstance().Generate())
}
