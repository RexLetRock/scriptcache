package zevio

import (
	"bufio"
	"bytes"

	"github.com/tidwall/evio"
)

const ACK_SEP = "|"
const ENDLINE = "#\t#"
const ENDLINE_LENGTH = len(ENDLINE)

func MainEvio(address string) {
	var events evio.Events
	events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		c.SetContext(&evio.InputStream{})
		// logrus.Warnf("New context %+v \n", c.RemoteAddr())
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

		msgs := bytes.Split(data, []byte(ENDLINE))
		if len(msgs) < 1 {
			return
		}

		msgsA := msgs[:len(msgs)-1]
		msgsB := msgs[len(msgs)-1]

		// Range data
		resdata := []byte{}
		for _, v := range msgsA {
			vf := bytes.Split(v, []byte(ACK_SEP))
			vfull := append(v, []byte(string(vf[0])+ENDLINE)...)
			resdata = append(resdata, vfull...)
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
