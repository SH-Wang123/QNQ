package QNQ

/*
*
其他handler向daemon注册信号后，会获得信号的值。handler根据该值响应信号
uint64
| 0   | 1 ~ 4 | 5 ~ 8 | 9 ~ 32 |
| 预留 | 类别   | 信号  | Task Id |
*/
const (
	typeBin   = 0x70000000
	signalBin = 0x07800000
	taskIdBin = 0x007fffff
)

// task type
const (
	_ = iota
)
