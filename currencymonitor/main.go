package main

import (
	"log"
	"net/http"

	"github.com/Dmk88/go_practice/currencymonitor/handlers"
	"github.com/Dmk88/go_practice/currencymonitor/monitor"
)

func main() {
	config := monitor.LoadConfig()
	log.Printf("%#v\n", config)

	scylla := monitor.InitScylla(config.Scylla)
	daemon := monitor.NewDaemon(scylla)
	// check after reboot
	daemon.CheckMonitoring()
	h := handlers.NewHandler(config, scylla, daemon)

	http.HandleFunc("/start", h.Start)
	http.HandleFunc("/results", h.MonitoringResults)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
