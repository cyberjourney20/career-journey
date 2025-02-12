package utils

import (
	"encoding/json"
	"log"
	"time"
)

func TimeLineBuilderJSON(t time.Time, s string) string {
	timeline := map[string]string{
		"marker": t.Format(time.RFC3339), // Use a standard format
		"event":  s,
	}
	// Marshal the map into JSON bytes
	jsonBytes, err := json.Marshal(timeline)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return ""
	}

	// Return the JSON as a string
	return string(jsonBytes)

	// ts := t.String()
	// json := fmt.Sprintf("{'marker':'%s', 'event':'%s'}", ts, s)
	// return json

}
