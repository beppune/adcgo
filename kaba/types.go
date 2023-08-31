package kaba

import "fmt"

//import "time"

type KabaEntry struct {
	TimeStamp string
	Support   string
	Name      string
	Sensor    string
	Text      string
}

func (k *KabaEntry) String() string {
	if k.Support == "" {
		return fmt.Sprintf("%v, %v, %v", k.TimeStamp, k.Sensor, k.Text)
	}
	return fmt.Sprintf("%v, %v, %v, %v, %v", k.TimeStamp, k.Support, k.Name, k.Sensor, k.Text)
}

func (k *KabaEntry) IsAllarm() bool {
	return k.Support == ""
}
