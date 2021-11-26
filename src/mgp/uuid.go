package mgp

import (
	"strconv"
	"time"
)

const UuidPrefix string = "ga"

var (
	machineID    int64 // range [ 0, 1023 ]
	sn            int64 // range [ 0, 4095 ]
	lastTimeStamp int64
)

func init() {
	lastTimeStamp = time.Now().UnixNano() / 1000000
}

func SetMachineId(mid int64) {
	machineID = mid << 12
}

func InitUuid() string  {
	r := genSnowflakeId()
	uuid := UuidPrefix + strconv.FormatInt(r, 10)

	return uuid
}

func genSnowflakeId() int64 {
	curTimeStamp := time.Now().UnixNano() / 1000000

	if curTimeStamp == lastTimeStamp {
		sn++
		// range [ 0, 4095 ]
		if sn > 4095 {
			time.Sleep(time.Millisecond)
			curTimeStamp = time.Now().UnixNano() / 1000000
			lastTimeStamp = curTimeStamp
			sn = 0
		}

		rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
		rightBinValue <<= 22
		id := rightBinValue | machineID | sn
		return id
	}

	if curTimeStamp > lastTimeStamp {
		sn = 0
		lastTimeStamp = curTimeStamp
		rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
		rightBinValue <<= 22
		id := rightBinValue | machineID | sn
		return id
	}

	if curTimeStamp < lastTimeStamp {
		return 0
	}

	return 0
}