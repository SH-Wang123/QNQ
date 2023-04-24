package common

import "log"

func init() {
	initLockMap()
	newFixedGoPool(0)
	if globalCoroPool == nil {
		globalCoroPool = NewFixedPool(0)
		log.Printf("Create a new coroutines pool, size : %v", globalCoroPool.GoNum)
	}
	GetCoroutinesPool().StartPool()
}
