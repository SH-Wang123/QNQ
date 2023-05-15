package network

func init() {
	go StartQServers()
	NetChan.StartPump()
}
