package tcp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"scriptcache/colf/message"
	"sync"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const ThreadPerConn = 5
const countSize = 5_00_000
const connStr = "developer:password@tcp(127.0.0.1:4000)/imsystem?parseTime=true"

var pCounter = PerformanceCounterCreate(countSize, 0, "SERVER RUN")
var GCache *DBCache

type DBCache struct {
	cache sync.Map
	db    *sql.DB
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

			// QUERY GROUP INFO
			lastMsgID, ok := GCache.cache.Load(m.GroupId)
			if !ok {
				// QUERY DATABASE
				var data uint64
				GCache.db.QueryRow(fmt.Sprintf("SELECT message_id FROM ims_message WHERE group_id=%v ORDER BY created_at DESC LIMIT 1", m.GroupId)).Scan(&data)
				fmt.Printf("%v \n", data)
				GCache.cache.Store(m.GroupId, data)
				lastMsgID = data
			} else {
				fmt.Printf("FROM CACHE %v \n", lastMsgID)
			}

			// select message_id from ims_message where group_id=381870481448962 order by created_at desc limit 1;
			// select * from ims_message where group_id=381870481448962 order by created_at desc limit 1;
			// INSERT INTO ims_message(message_id, group_id, data, flags, created_at) VALUES(?, ?, ?, 0, ?)
			// developer:password@tcp(127.0.0.1:4000)/imsystem?parseTime=true

			// SEND
			reMsg := append(msg[0:4], []byte(fmt.Sprintf("%v#\t#", binary.LittleEndian.Uint32(msg[0:4])))...)
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
