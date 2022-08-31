package tcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"scriptcache/colf/message"
	"sync"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const ThreadPerConn = 5
const countSize = 5_00_000
const connStr = "developer:password@tcp(127.0.0.1:4000)/imsystem?parseTime=true"
const cDebug = true
const cSnowflakeNode = 1
const createMessageStmt = "INSERT INTO ims_message(message_id, group_id, data, flags, created_at) VALUES(?, ?, ?, 0, ?)"

var pCounter = PerformanceCounterCreate(countSize, 0, "SERVER RUN")
var GCache *DBCache

type DBCache struct {
	cache sync.Map
	db    *sql.DB
	nf    *Snowflake
}

func DBCacheCreate(dburl string) *DBCache {
	db, err := sql.Open("mysql", dburl)
	if err != nil {
		log.Fatalf("Cant connect db", err)
	}

	s := &DBCache{
		db:    db,
		cache: sync.Map{},
	}
	s.nf, _ = NewSnowflake(cSnowflakeNode)

	return s
}

func ServerStart() {
	GCache = DBCacheCreate(connStr)

	listener, _ := net.Listen("tcp", "0.0.0.0:8888")
	defer listener.Close()
	for {
		if conn, err := listener.Accept(); err == nil {
			go handleConn(conn)
		}
	}
}

func handleConn(conn net.Conn) {
	handle := ConnHandleCreate(conn)
	handle.Handle()
}

type ConnHandle struct {
	readerReq *io.PipeReader
	writerReq *io.PipeWriter

	buffer []byte
	slice  []byte
	conn   net.Conn
}

func showLog(format string, v ...any) {
	if cDebug {
		log.Printf(format, v...)
	}
}

func ConnHandleCreate(conn net.Conn) *ConnHandle {
	p := &ConnHandle{
		buffer: make([]byte, 1024*1000),
		conn:   conn,
	}
	p.readerReq, p.writerReq = io.Pipe()

	// Request flow
	go func() {
		reader := bufio.NewReader(p.readerReq)
		cSend := 0
		for {
			msg, _ := readWithEnd(reader)

			// DECODE
			m := message.Message{}
			m.Unmarshal(msg[4 : len(msg)-3])

			// HANDLE
			// Query lastest messageid info
			lastMsgID, ok := GCache.cache.Load(m.GroupId)
			if !ok {
				var data uint64
				GCache.db.QueryRow(fmt.Sprintf("SELECT message_id FROM ims_message WHERE group_id=%v ORDER BY created_at DESC LIMIT 1", m.GroupId)).Scan(&data)
				if data != 0 {
					showLog("db %v \n", data)
					GCache.cache.Store(m.GroupId, data)
					lastMsgID = data
				}
			} else {
				showLog("cached %v \n", lastMsgID)
			}

			nextMsgID := uint64(0)
			if lastMsgID != nil {
				nextMsgID = GCache.nf.NextIdWithSeq(lastMsgID.(uint64))
				GCache.cache.Store(m.GroupId, nextMsgID)

				// Save to cache for flush to DB
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
					defer cancel()

					stmt, err := GCache.db.PrepareContext(ctx, createMessageStmt)
					if err != nil {
						showLog("Sql error %v \n", err)
					}
					defer stmt.Close()

					if _, err = stmt.ExecContext(ctx, nextMsgID, m.GroupId, m.Data, time.Now().UnixMicro()+(7*3600)); err != nil {
						showLog("Sql error %v \n", err)
					}
				}()
			}

			// SEND
			reMsg := append(msg[0:4], []byte(fmt.Sprintf("%v#\t#", nextMsgID))...)
			cSend += 1
			if true {
				p.conn.Write(reMsg)
			} else {
				p.slice = append(p.slice, reMsg...)
				if cSend%cSendSize == 0 {
					p.conn.Write(p.slice)
					p.slice = []byte{}
				}
			}
			pCounter.Step(true)
		}
	}()

	// Response flow
	return p
}

func (s *ConnHandle) Handle() error {
	defer s.conn.Close()
	for {
		n, err := s.conn.Read(s.buffer)
		if err != nil {
			return err
		}
		s.writerReq.Write(s.buffer[:n])
	}
}
