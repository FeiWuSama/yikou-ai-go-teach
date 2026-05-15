package snowflake

import (
	"github.com/sony/sonyflake"
	"strconv"
)

var (
	sf *sonyflake.Sonyflake
)

func init() {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (uint16, error) { return 1, nil },
	})
}

func GenerateSnowFlakeId() (int64, error) {
	id, err := sf.NextID()
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func GenerateSnowFlakeIdString() (string, error) {
	snowFlakeId, err := GenerateSnowFlakeId()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(snowFlakeId)), nil
}
