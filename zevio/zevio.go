package zevio

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"github.com/RexLetRock/scriptcache/ztcp/ztcputil"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/evio"
)

const FRAMESPLIT = "|"
const ENDLINE = "#\t#"
const ENDLINE_LENGTH = len(ENDLINE)

var ENDBYTE = []byte(ENDLINE)

var DataGroup ztcputil.ConcurrentMap // [CMaxResultBuffer]*[]byte
var DataIP ztcputil.ConcurrentMap

func init() {
	DataGroup = ztcputil.CMapCreate()
	DataIP = ztcputil.CMapCreate()
}

func MainEvio(address string) {
	var events evio.Events
	events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		c.SetContext(&evio.InputStream{})
		logrus.Warnf("New context %+v \n", c.RemoteAddr())
		return
	}

	events.Closed = func(c evio.Conn, err error) (action evio.Action) {
		return
	}

	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		if in == nil {
			return
		}

		is := c.Context().(*evio.InputStream)
		data := is.Begin(in)

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
				}
			}

			resmsg := append([]byte(vresp), ENDBYTE...)
			resdata = append(resdata, resmsg...)
		}

		// Leftover
		is.End(msgsB)
		out = resdata
		return
	}

	if err := evio.Serve(events, address); err != nil {
		panic(err.Error())
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
