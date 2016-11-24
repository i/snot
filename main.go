package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	readDevice, err := os.Open("/dev/input/event0")
	f, err := os.OpenFile("/dev/input/event0", os.O_APPEND|os.O_WRONLY, os.ModeCharDevice|os.ModeAppend)
	if err != nil {
		//		log.Fatalf("Error opening path: %v", err)
	}
	fmt.Println("OK")
	fmt.Println(f)

	//mappings := make(map[string]map[string][]byte)
	var mappings map[string]map[string][]byte
	/*
		for _, letter := range characters {
			for i:=0; i < 2; i++ {
				buf := make([]byte, 128)
				fmt.Printf("type character %q\n", letter)
				n, err := readDevice.Read(buf)
				if err != nil {
					log.Fatalf("error reading: %v", err)
				}
				fmt.Printf("read %d bytes\n", n)
				mappings[strings.ToUpper(letter)] = append(mappings[strings.ToUpper(letter)],buf[:n])
			}
		}

		jsn, err := json.Marshal(mappings)
		if err != nil {
			log.Fatalf("error writing mappings: %v", err)
		}
		fmt.Println(string(jsn))

		mapFile, err := os.Create("./mappings_v2.json")
		if err != nil {
			log.Fatalf("error creating mapping file: %v", err)
		}
		mapFile.Write(jsn)
		mapFile.Close()
		os.Exit(0)
	*/

	mapFile, err := os.Open("./mappings_v2.json")
	if err != nil {
		log.Fatalf("error reading mappings file: %v", err)
	}
	if err := json.NewDecoder(mapFile).Decode(&mappings); err != nil {
		log.Fatalf("error decoding mappings file: %v", err)
	}
	fmt.Println(readDevice)

	fmt.Println(mappings)

	// writeChar := func(char string) {
	// 	if codes, ok := mappings[char]; ok {
	// 		for _, code := range codes {
	// 			if _, err := f.Write(code); err != nil {
	// 				log.Fatal("error writing to device: %v", err)
	// 			}
	// 		time.Sleep(time.Millisecond*20)
	// 		}
	// 	}
	// }

	charDown := func(char string) {
		fmt.Println("charDown:", char)
		if codes, ok := mappings[char]; ok {
			if _, err := f.Write(codes["down"]); err != nil {
				log.Fatal("error writing to device: %v", err)
			}
		}
	}

	charUp := func(char string) {
		fmt.Println("charUp:", char)
		if codes, ok := mappings[char]; ok {
			if _, err := f.Write(codes["up"]); err != nil {
				log.Fatal("error writing to device: %v", err)
			}
		}
	}

	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		for {
			msg := make([]byte, 32)
			n, err := ws.Read(msg)
			if err != nil {
				log.Println(err)
				return
			}
			parts := strings.Split(string(msg[:n]), ".")
			if len(parts) != 2 {
				fmt.Printf("invalid args: %v\n", parts)
				continue
			}
			switch parts[1] {
			case "up":
				charUp(strings.ToUpper(parts[0]))
			case "down":
				charDown(strings.ToUpper(parts[0]))
			default:
				fmt.Println("invalid op: %v", parts[0])
			}
			fmt.Println(string(msg))
		}
	}))
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "./public/index.html")
	//})
	http.Handle("/", http.FileServer(http.Dir("public")))

	go http.ListenAndServe(":3000", nil)

	input := make([]byte, 128)
	for {
		n, err := os.Stdin.Read(input)
		if err != nil {
			log.Fatalf("error reading stdin: %v", err)
		}
		charDown(strings.ToUpper(strings.TrimSpace(string(input[:n]))))
		time.Sleep(time.Millisecond * 20)
		charUp(strings.ToUpper(strings.TrimSpace(string(input[:n]))))
	}
}

// if _, err := f.Write(buf[:n]); err != nil {
// 	fmt.Printf("Couldn't write to file: %v", err.Error())
// }

var characters = []string{
	// 	"A", "B", "C", "D", "E", "F", "G", "H",
	// 	"I", "J", "K", "L", "M", "N", "O", "P",
	// 	"Q", "R", "S", "T", "U", "V", "W", "X",
	//	"Y", "Z", "LEFT", "RIGHT", "ESC",
	// "UP", "DOWN",
	"C",
}

// for each character
//   prompt user for input
//   record input bytes
//   set map[input] to bytes
// for each event
// send map[input] bytes to device
