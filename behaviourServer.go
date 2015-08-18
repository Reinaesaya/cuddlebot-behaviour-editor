package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type SetPIDParams struct {
	Addr string
	Kp   float32
	Ki   float32
	Kd   float32
}

type SetPointParams struct {
	Addr      string
	Delay     uint16
	Loop      uint16
	Setpoints []uint16
}

type SleepParams struct {
	Addr []string
}

type SmoothParams struct {
	Addr		string
	Time		uint16
	Setpoint	[]uint16
}

type BehaviourParams struct {
	FileName           string
	BehaviourName      string
	Data               string
	OverwriteFile      bool
	OverwriteBehaviour bool
}

type GestureParams struct {
	Name string
}

var behaviourNameToDataMap map[string]string

var DEFAULT_PATH string

var currGesture = "unknown"

func sendCommand(commandBytes []byte, commandstr string) {
	var url string
	switch (commandstr) {
	case "setpid":
		url = "http://10.10.10.1/1/setpid.json"
	case "setpoint":
		url = "http://10.10.10.1/1/setpoint.json"
	case "smooth":
		url = "http://10.10.10.1/1/smooth.json"
	case "sleep":
		url = "http://10.10.10.1/1/sleep.json"
	default:
		log.Println("Not a valid command")
		return
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(commandBytes))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	} else {
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
}

func mainView(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
	return
}

func setpid(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	decoder := json.NewDecoder(req.Body)
	var setpid SetPIDParams
	err := decoder.Decode(&setpid)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(setpid)

		var setpidParams SetPIDParams
		setpidParams.Addr = setpid.Addr
		setpidParams.Kp = setpid.Kp
		setpidParams.Ki = setpid.Ki
		setpidParams.Kd = setpid.Kd
		jsonBytes, err := json.Marshal(setpidParams)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			go sendCommand(jsonBytes, "setpid")
		}
	}
}

func setpoint(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	decoder := json.NewDecoder(req.Body)
	var setPt SetPointParams
	err := decoder.Decode(&setPt)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(setPt)

		var setPtParams SetPointParams
		setPtParams.Addr = setPt.Addr
		setPtParams.Delay = setPt.Delay
		setPtParams.Loop = setPt.Loop
		setPtParams.Setpoints = setPt.Setpoints
		jsonBytes, err := json.Marshal(setPtParams)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			go sendCommand(jsonBytes, "setpoint")
		}
	}
}

func smooth(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	decoder := json.NewDecoder(req.Body)
	var smooths SmoothParams
	err := decoder.Decode(&smooths)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(smooths)

		var smoothinstructs SmoothParams
		smoothinstructs.Addr = smooths.Addr
		smoothinstructs.Time = smooths.Time
		smoothinstructs.Setpoint = smooths.Setpoint
		jsonBytes, err := json.Marshal(smoothinstructs)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			go sendCommand(jsonBytes, "smooth")
		}
	}
}

func sleep(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	decoder := json.NewDecoder(req.Body)
	var sleeps SleepParams
	err := decoder.Decode(&sleeps)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(sleeps)
		log.Println(sleeps.Addr)
		
		var sleepParams SleepParams
		// all := []string{"all"}
		// if sleeps.Addr == all {
		//	sleepParams := SleepParams{
		//		Addr: []string{"purr", "headx", "heady", "ribs", "spine"}, }
		// } else {
			sleepParams.Addr = sleeps.Addr
		//}
		jsonBytes, err := json.Marshal(sleepParams)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			go sendCommand(jsonBytes, "sleep")
		}
	}
}

func saveBehaviourParams(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// if origin := req.Header.Get("Origin"); origin != "" {
	//  rw.Header().Set("Access-Control-Allow-Origin", origin)
	// }

	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	decoder := json.NewDecoder(req.Body)
	var behParams BehaviourParams
	err := decoder.Decode(&behParams)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(behParams)
		log.Println("BehaviourParams to save: ", behParams.BehaviourName)

		if _, err := os.Stat(behParams.FileName); err == nil && !behParams.OverwriteFile {
			log.Printf("File \"%s\" already exists, verifying overwrite", behParams.BehaviourName)
			rw.Write([]byte("overwriteFile"))
		} else if _, ok := behaviourNameToDataMap[behParams.BehaviourName]; ok && !behParams.OverwriteBehaviour {
			log.Printf("Behaviour \"%s\" already exists, verifying overwrite", behParams.BehaviourName)
			rw.Write([]byte("overwriteBehaviour"))
		} else {
			behaviourNameToDataMap[behParams.BehaviourName] = behParams.Data

			writeBehavioursToFile(behParams.FileName)
		}
	}
}

func loadBehaviourParams(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// if origin := req.Header.Get("Origin"); origin != "" {
	//  rw.Header().Set("Access-Control-Allow-Origin", origin)
	// }

	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	bytes, err := ioutil.ReadAll(req.Body)

	var path string

	if err != nil {
		log.Println(err)
	} else {
		if string(bytes) == "defaults" {
			path = DEFAULT_PATH
		} else {
			path = "./" + string(bytes)
		}

		log.Printf("loadBehaviourParams: path: %v\n", path)
		behsStr, err := loadBehavioursFromFile(path)

		if err != nil {
			log.Println(err)
		} else {

			rw.Write([]byte(behsStr))

			var data map[string]interface{}

			if err := json.Unmarshal(bytes, &data); err != nil {
				log.Printf("Unable to load behaviours from file: %v\n", path)
				log.Println(err)
			} else {

				for name, data := range data {
					log.Printf("name: %v, data: %v\n", name, data)
				}
			}
		}
	}
}

func gesture(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// if origin := req.Header.Get("Origin"); origin != "" {
	//  rw.Header().Set("Access-Control-Allow-Origin", origin)
	// }

	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	rw.Write([]byte(currGesture))
}

func writeBehavioursToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	fmt.Fprint(w, "[")
	started := false
	for name, data := range behaviourNameToDataMap {
		if started {
			fmt.Fprint(w, ",")
		}
		started = true
		fmt.Fprint(w, "{\"Name\":\""+name+"\", \"Data\":"+data+"}")
	}
	fmt.Fprint(w, "]")

	return w.Flush()
}

func loadBehavioursFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	allBehsString := ""
	for scanner.Scan() {
		var line = scanner.Text()
		//      log.Println("line: ", line)
		allBehsString += line
	}

	return allBehsString, scanner.Err()
}

func listenForGestureCommands(ipPortStr string) {

	udpAddr, err := net.ResolveUDPAddr("udp", ipPortStr)
	if err != nil {
		log.Println(err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println(err)
		return
	}

	for {

		fmt.Println("reading udp")

		readData := make([]byte, 1000)
		_, _, err := udpConn.ReadFromUDP(readData)

		fmt.Println("read from udp")

		if err != nil {
			log.Println(err)
			currGesture = "unknown"
		} else {
			msg := string(readData)
			currGesture = msg
		}
	}
}

func main() {
	DEFAULT_PATH = "./DefaultBehaviours.txt"

	//  loadBehavioursFromFile(DEFAULT_PATH)

	behaviourNameToDataMap = make(map[string]string)

	ipPortStr := ":1234"
	go listenForGestureCommands(ipPortStr)

	http.HandleFunc("/gesture", gesture)
	http.HandleFunc("/setpid", setpid)
	http.HandleFunc("/setpoint", setpoint)
	http.HandleFunc("/smooth", smooth)
	http.HandleFunc("/sleep", sleep)
	http.HandleFunc("/savebehaviour", saveBehaviourParams)
	http.HandleFunc("/loadbehaviour", loadBehaviourParams)

	http.Handle("/", http.FileServer(http.Dir("./web")))

	//  http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/assets"))))
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("ListenAndServer err: %s\n", err)
	}
}
