package tcp

import (
	"sync"
	"time"
)

var (
	kMaxUInt64          uint64 = 0xFFFFFFFFFFFFFFFF
	kEpoch              uint64 = 1288834974657
	kWorkerIdBits       uint64 = 14
	kWorkerMaxId        uint64 = kMaxUInt64 ^ (kMaxUInt64 << kWorkerIdBits)
	kWorkerIdMask       uint64 = kMaxUInt64 ^ (kMaxUInt64 << kWorkerIdBits)
	kSequenceBits       uint64 = 12
	kSequenceMask       uint64 = kMaxUInt64 ^ (kMaxUInt64 << kSequenceBits)
	kWorkerIdShift      uint64 = kSequenceBits
	kTimestampLeftShift uint64 = kSequenceBits + kWorkerIdBits
)

type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp uint64
	workerId      uint64
	seqId         uint64
	lastSeqIds    map[uint64]bool
}

func miliseconds() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}

func nanoseconds() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond)) // int64(1_000_000))
}

func snowtime() uint64 {
	return miliseconds()
}

func NewSnowflake(workerId uint64) (*Snowflake, error) {
	if workerId > kWorkerMaxId {
		workerId = workerId & kWorkerIdMask
	}
	snowflake := &Snowflake{
		lastTimestamp: kMaxUInt64,
		workerId:      workerId,
		seqId:         uint64(0),
		lastSeqIds:    make(map[uint64]bool),
	}
	return snowflake, nil
}

func (s *Snowflake) NextId() uint64 {
	s.mu.Lock()
	now := snowtime()
	if now == s.lastTimestamp {
		s.seqId = (s.seqId + 1) & kSequenceMask
		if s.seqId == 0 {
			//wait until next milisecond
			for now <= s.lastTimestamp {
				now = snowtime()
			}
		}
	} else {
		s.seqId = 0
	}
	s.lastTimestamp = now

	r := (now-kEpoch)<<kTimestampLeftShift |
		(s.workerId << kWorkerIdShift) |
		s.seqId

	s.mu.Unlock()
	return r
}

func (s *Snowflake) NextIdWithSeq(lastSeqId uint64) uint64 {
	s.mu.Lock()
	now := snowtime()
	lastSeqId = (lastSeqId + 1) & kSequenceMask
	if now == s.lastTimestamp {
		if _, found := s.lastSeqIds[lastSeqId]; found {
			//wait until next milisecond
			for now <= s.lastTimestamp {
				now = snowtime()
			}
			s.lastSeqIds = make(map[uint64]bool)
		} else {
			s.lastSeqIds[lastSeqId] = true
		}
	} else {
		s.lastSeqIds = make(map[uint64]bool)
	}
	s.lastTimestamp = now

	r := (now-kEpoch)<<kTimestampLeftShift |
		(s.workerId << kWorkerIdShift) |
		lastSeqId

	s.mu.Unlock()
	return r
}
