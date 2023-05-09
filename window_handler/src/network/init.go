package network

func init() {
	go StartQTargets()
	go StartQServers()
	NetChan.StartPump()
}
