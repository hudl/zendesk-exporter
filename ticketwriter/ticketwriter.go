package ticketwriter

import (
	"github.com/adamar/ZeGo/zego"
	"github.com/pivotal-golang/bytefmt"

	"encoding/json"
	"os"
)

func Write(ticks []zego.Ticket) {
	tickFile, _ := os.OpenFile("./tickets.json", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	for _, t := range ticks {
		jsonBytes, err := json.Marshal(t)
		if err != nil {
			log.Error("Error marshalling ticket to json. Ticket: %+v \nError: %+v", t, err)
			continue
		}
		_, err = tickFile.WriteString(string(jsonBytes))
		if err != nil {
			log.Error("Error writing to file: %+v", err)
		}
		tickFile.WriteString("\n")
	}
	st, _ := tickFile.Stat()
	szStr := bytefmt.ByteSize(uint64(st.Size()))
	log.Info("File is %s large", szStr)
}
