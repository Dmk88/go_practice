package monitor

import (
	"fmt"
	"math/rand"
	"time"
)

func genRsp() string {
	rand.Seed(time.Now().UnixNano())
	min := 39000
	max := 41000
	amount := rand.Intn(max-min) + min
	return fmt.Sprintf(`{ "amount": %d }`, amount)
}
