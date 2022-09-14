package ztcputil

import "time"

const NCpu = 60
const NRun = 30_000_000
const ENDLINE = "#\t#"

const SendSize = 10_000
const ChansSize = 1024 * 100
const TimeToFlush = 10 * time.Microsecond
