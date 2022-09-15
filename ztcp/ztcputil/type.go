package ztcputil

import "time"

const NCpu = 12
const NRun = 50_000_000
const ENDLINE = "#\t#"

const SendSize = 10_000
const ChansSize = 1024 * 100
const TimeToFlush = 1000 * time.Microsecond
