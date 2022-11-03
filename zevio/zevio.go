package zevio

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"github.com/RexLetRock/scriptcache/ztcp/ztcputil"
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
var BroadcastCount ztcputil.Count32

type EvioContext struct {
	// is *evio.InputStream
}

func init() {
	DataGroup = ztcputil.CMapCreate()
	DataIP = ztcputil.CMapCreate()
	ChanBroadcast = make(chan []byte, 1000)
	go Broadcast()
}

func MainEvio(address string) {
	var events evio.Events
	events.NumLoops = ztcputil.NCpu
	events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		c.SetContext(&evio.InputStream{})
		DataIP.Set(c.RemoteAddr().String(), &c)
		return
	}

	events.Closed = func(c evio.Conn, err error) (action evio.Action) {
		return
	}

	events.Cast = func(c evio.Conn) (out []byte) {
		out = []byte("hello" + ENDLINE)
		return
	}

	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		ctx := c.Context().(*evio.InputStream)
		if in == nil {
			// return
			in = []byte{}
		}

		data := ctx.Begin(in)
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
						groupMessID = ztcputil.Count32Create()
					}
					groupMessIDInt, _ := groupMessID.(*ztcputil.Count32)
					groupMessIDInt.Inc()
					DataGroup.Set(vdata[2], groupMessIDInt)
					vresp += FRAMESPLIT + strconv.Itoa(int(groupMessIDInt.Get()))
				case MessageBroadcast.Toa():
					ChanBroadcast <- []byte{}
				}
			}

			resmsg := append([]byte(vresp), ENDBYTE...)
			resdata = append(resdata, resmsg...)
		}

		// Leftover
		ctx.End(msgsB)
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
			con := v.(*evio.Conn)
			if con != nil {

			}
		}
	}
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
