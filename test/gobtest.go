package main

import "github.com/bennicholls/burl-E/burl"
import "os"
import "fmt"
import "encoding/gob"

func LogStat(s burl.Stat) {
	fmt.Println("min: ", s.Min(), " val: ", s.Get(), " max: ", s.Max())
}

//to test composite structs
type Stats struct {
	A burl.Stat
	B burl.Stat
	Text string
}

func main() {

	f, err := os.Create("save.txt")
	if err != nil {
		fmt.Println("cant make save file: " + err.Error())
		return
	}
	
	ToSave := Stats{A: burl.NewStat(50), B: burl.NewStat(90), Text: "whatever"}
	ToSave2 := Stats{A: burl.NewStat(40), B: burl.NewStat(70), Text: "whatever2"}
	var ToLoad Stats
	var ToLoad2 Stats

	ToSave.A.Mod(-10)
	ToSave.B.SetMin(30)
	LogStat(ToSave.A)
	LogStat(ToSave.B)
	fmt.Println(ToSave.Text)

	LogStat(ToSave2.A)
	LogStat(ToSave2.B)
	fmt.Println(ToSave2.Text)

	enc := gob.NewEncoder(f)
	err = enc.Encode(ToSave)
	err = enc.Encode(ToSave2)
	if err != nil {
		fmt.Println("Could not encode: " + err.Error())
	}
	f.Close()

	load, err := os.Open("save.txt")
	if err != nil {
		fmt.Println("cant load save file: " + err.Error())
		return
	}

	dec := gob.NewDecoder(load)
	err = dec.Decode(&ToLoad)
	err = dec.Decode(&ToLoad2)
	if err != nil {
		fmt.Println("Could not decode: " + err.Error())
	}

	LogStat(ToLoad.A)
	LogStat(ToLoad.B)
	fmt.Println(ToLoad.Text)

	LogStat(ToLoad2.A)
	LogStat(ToLoad2.B)
	fmt.Println(ToLoad2.Text)
	

	load.Close()
}