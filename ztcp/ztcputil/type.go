package ztcputil

import "time"

const NCpu = 12
const NRun = 10_000_000
const FRAMESPLIT = "|"
const ENDLINE = "#\t#" //
const ENDLINE_LENGTH = len(ENDLINE)

const SendSize = 999
const ChanSize = 1000
const BuffSize = 1000
const CRound = 12
const TimeToFlush = time.Millisecond

var ENDBYTE = []byte(ENDLINE)
