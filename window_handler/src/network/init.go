package network

func init() {
	go StartQClient()
	go StartQServer()
	NetChan.StartPump()
}
