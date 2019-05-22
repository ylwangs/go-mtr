package common

import (
	"time"
)

const DEFAULT_PORT = 33434
const DEFAULT_MAX_HOPS = 30
const DEFAULT_TIMEOUT_MS = 800

const DEFAULT_RETRIES = 5
const DEFAULT_PACKET_SIZE = 56
const DEFAULT_SNT_SIZE = 10

type IcmpReturn struct {
	Success bool
	Addr    string
	Elapsed time.Duration
}

type IcmpHop struct {
	Success  bool
	Address  string
	Host     string
	N        int
	TTL      int
	Snt      int
	LastTime time.Duration
	AvgTime  time.Duration
	BestTime time.Duration
	WrstTime time.Duration
	Loss     float32
}

type MtrReturn struct {
	Success  bool
	TTL      int
	Host     string
	SuccSum  int
	LastTime time.Duration
	AllTime  time.Duration
	BestTime time.Duration
	AvgTime  time.Duration
	WrstTime time.Duration
	Seq      int
}

type MtrResult struct {
	DestAddress string
	Hops        []IcmpHop
}

type MtrOptions struct {
	port       int
	maxHops    int
	timeoutMs  int
	retries    int
	packetSize int
	sntSize    int
}

func (options *MtrOptions) Port() int {
	if options.port == 0 {
		options.port = DEFAULT_PORT
	}
	return options.port
}

func (options *MtrOptions) SetPort(port int) {
	options.port = port
}

func (options *MtrOptions) MaxHops() int {
	if options.maxHops == 0 {
		options.maxHops = DEFAULT_MAX_HOPS
	}
	return options.maxHops
}

func (options *MtrOptions) SetMaxHops(maxHops int) {
	options.maxHops = maxHops
}

func (options *MtrOptions) TimeoutMs() int {
	if options.timeoutMs == 0 {
		options.timeoutMs = DEFAULT_TIMEOUT_MS
	}
	return options.timeoutMs
}

func (options *MtrOptions) SetTimeoutMs(timeoutMs int) {
	options.timeoutMs = timeoutMs
}

func (options *MtrOptions) Retries() int {
	if options.retries == 0 {
		options.retries = DEFAULT_RETRIES
	}
	return options.retries
}

func (options *MtrOptions) SetRetries(retries int) {
	options.retries = retries
}

func (options *MtrOptions) SntSize() int {
	if options.sntSize == 0 {
		options.sntSize = DEFAULT_SNT_SIZE
	}
	return options.sntSize
}

func (options *MtrOptions) SetSntSize(sntSize int) {
	options.sntSize = sntSize
}

func (options *MtrOptions) PacketSize() int {
	if options.packetSize == 0 {
		options.packetSize = DEFAULT_PACKET_SIZE
	}
	return options.packetSize
}

func (options *MtrOptions) SetPacketSize(packetSize int) {
	options.packetSize = packetSize
}
