package ztcp

import "time"

const NCpu = 12
const NRun = 30_000_000
const ENDLINE = "#\t#"

const cSendSize = 10_000
const cChansSize = 1024 * 100
const cTimeToFlush = 1000 * time.Microsecond
