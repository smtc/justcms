package database

import (
	"net"
	"os"

	"github.com/guotie/config"
	"github.com/guotie/deferinit"
)

var (
	securekey []byte
)

func init() {
	deferinit.AddInit(OpenDefaultDB, nil, 1000)
	//deferinit.AddInit(openRedis, CloseRedis, 1100)

	deferinit.AddInit(func() {
		securekey = config.GetBytesDefault("securekey", []byte("eTRp0ae9wg5hIyU3Re2UzOungcF6raDjYUXGTfsFxm6jH41FlGpWO7q5ndP9UHO"))
		// pid
		pid = _genpid()

		// machineid 向左移位20
		machineid = _genmachineid() << 20
	}, nil, 0)
}

// 20 bit
func _genpid() uint64 {
	id := os.Getpid()
	id = ((id & 0xfffff) + (id >> 20)) & 0xfffff

	return uint64(id)
}

// 24位
func _genmachineid() uint64 {
	var s string

	ifs, err := net.Interfaces()
	if err == nil {
		for _, i := range ifs {
			s += i.HardwareAddr.String()
		}
	} else {
		s = err.Error()
	}

	code := hashcode(s)
	code = ((code & 0xffffff) + (code >> 24)) & 0xffffff

	return uint64(code)
}

func hashcode(s string) uint32 {
	var val uint32 = 1

	for i := 0; i < len(s); i++ {
		val += (val * 37) + uint32(s[i])
	}

	return val
}
