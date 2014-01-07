package survana

import (
	"math/rand"
	"time"
)

const (
	ID_STRING = "AaBbCcDdEeFfGgHhJjKLMmNnoPpqRrSsTtUuVvWwXxYyZz3456789"
)

var (
	NID_STRING int        = len(ID_STRING)
	generator  *rand.Rand = nil
)

func RandomId(length int) (id string) {
	for i := 0; i < length; i++ {
		//append a random character from ID_STRING
		id += string(ID_STRING[generator.Intn(NID_STRING)])
	}

	return
}

func init() {
	//seed the PRNG source
	rand_src := rand.NewSource(time.Now().UnixNano())
	generator = rand.New(rand_src)
}
