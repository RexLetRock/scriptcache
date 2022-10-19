package zevio

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
	"sync"

	"github.com/RexLetRock/scriptcache/ztcp/ztcputil"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/evio"
)

const FRAMESPLIT = "|"
const ENDLINE = "#\t#"
const ENDLINE_LENGTH = len(ENDLINE)
const CHANSIZE = 10 * 1000

var ENDBYTE = []byte(ENDLINE)

var DataGroup ztcputil.ConcurrentMap // [CMaxResultBuffer]*[]byte
var DataIP ztcputil.ConcurrentMap
var ChanBroadcast chan []byte

type EvioStruct struct {
	is        *evio.InputStream
	ch        []byte
	mu        sync.Mutex
	testValue int
}

func init() {
	ChanBroadcast = make(chan []byte, 1000)
	go Broadcast()

	DataGroup = ztcputil.CMapCreate()
	DataIP = ztcputil.CMapCreate()
}

func MainEvio(address string) {
	var events evio.Events
	events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		evioStruct := &EvioStruct{
			is: &evio.InputStream{},
		}
		c.SetContext(evioStruct)
		DataIP.Set(c.RemoteAddr().String(), c)
		return
	}

	events.Closed = func(c evio.Conn, err error) (action evio.Action) {
		return
	}

	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		if in == nil {
			ctx := c.Context().(*EvioStruct)
			ctx.mu.Lock()
			if len(ctx.ch) > 0 {
				out = ctx.ch
				count := len(strings.Split(string(ctx.ch), ENDLINE))
				logrus.Warn("NILL ", string(out), count)
				ctx.ch = []byte{}
			}
			ctx.mu.Unlock()
			return
		}

		ctx := c.Context().(*EvioStruct)
		data := ctx.is.Begin(in)

		if len(data) < ENDLINE_LENGTH {
			return
		}

		msgs := bytes.Split(data, ENDBYTE)
		if len(msgs) < 1 {
			return
		}

		msgsA := msgs[:len(msgs)-1]
		msgsB := msgs[len(msgs)-1]

		// Range data
		resdata := []byte{}
		for _, v := range msgsA {
			vdata := strings.Split(string(v), FRAMESPLIT)
			vresp := vdata[0]
			if len(vdata) >= 3 {
				switch vdata[1] {
				case MessageNew.Toa():
					groupMessID, _ := DataGroup.Get(vdata[2])
					if groupMessID == nil {
						groupMessID = 0
					}
					groupMessIDInt, _ := groupMessID.(int)
					groupMessIDInt++
					DataGroup.Set(vdata[2], groupMessIDInt)
					vresp += FRAMESPLIT + strconv.Itoa(groupMessIDInt)
					ctx.testValue += 1
				case MessageBroadcast.Toa():
					ChanBroadcast <- []byte{}
				}
			}

			resmsg := append([]byte(vresp), ENDBYTE...)
			resdata = append(resdata, resmsg...)
		}

		// Leftover
		ctx.is.End(msgsB)
		ctx.mu.Lock()
		if len(ctx.ch) > 0 {
			logrus.Warn("DATA ", string(ctx.ch))
		}
		ctx.mu.Unlock()
		out = resdata
		return
	}

	if err := evio.Serve(events, address); err != nil {
		panic(err.Error())
	}
}

func Broadcast() {
	for range ChanBroadcast {
		ipdata := DataIP.Items()
		for _, v := range ipdata {
			con := v.(evio.Conn)
			ctx := con.Context().(*EvioStruct)
			ctx.mu.Lock()
			ctx.ch = append(ctx.ch, append([]byte("Hello"), ENDBYTE...)...)
			ctx.mu.Unlock()
			// con.Wake()
		}
	}
}

func CountTest() {
	ipdata := DataIP.Items()
	total := 0
	for _, v := range ipdata {
		total += v.(evio.Conn).Context().(*EvioStruct).testValue
	}
	logrus.Warn("TOTAL ", total)
}

func ReadWithEnd(reader *bufio.Reader) ([]byte, error) {
	message, err := reader.ReadBytes('#')
	if err != nil {
		return nil, err
	}

	a1, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	message = append(message, a1)
	if a1 != '\t' {
		message2, err := ReadWithEnd(reader)
		if err != nil {
			return nil, err
		}
		ret := append(message, message2...)
		return ret, nil
	}

	a2, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	message = append(message, a2)
	if a2 != '#' {
		message2, err := ReadWithEnd(reader)
		if err != nil {
			return nil, err
		}
		ret := append(message, message2...)
		return ret, nil
	}
	return message, nil
}
