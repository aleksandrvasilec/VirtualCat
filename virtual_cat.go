// virtual_cat.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

const (
	reset   = "\033[0m"
	green   = "\033[92m"
	red     = "\033[91m"
	yellow  = "\033[93m"
	blue    = "\033[94m"
	magenta = "\033[95m"
)

func colorize(text, color string) string {
	return color + text + reset
}

const catASCII = `
 /\_/\
( o.o )
 > ^ <
`

type State struct {
	Name       string `json:"name"`
	Hunger     int    `json:"hunger"`
	Happiness  int    `json:"happiness"`
	Energy     int    `json:"energy"`
	Health     int    `json:"health"`
	Age        int    `json:"age"`
	LastUpdate string `json:"last_update"`
}

func getHomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

func getConfigFile() string {
	return filepath.Join(getHomeDir(), ".virtual_cat.json")
}

func loadState() State {
	var state State
	data, err := ioutil.ReadFile(getConfigFile())
	if err != nil {
		return State{
			Name:       "Барсик",
			Hunger:     50,
			Happiness:  70,
			Energy:     80,
			Health:     90,
			Age:        0,
			LastUpdate: time.Now().Format(time.RFC3339),
		}
	}
	json.Unmarshal(data, &state)
	return state
}

func saveState(state State) {
	data, _ := json.MarshalIndent(state, "", "  ")
	ioutil.WriteFile(getConfigFile(), data, 0644)
}

func nowISO() string {
	return time.Now().Format(time.RFC3339)
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func getEmotion(state State) string {
	if state.Happiness >= 80 {
		return "😺"
	} else if state.Happiness >= 50 {
		return "😸"
	} else if state.Happiness >= 30 {
		return "😼"
	}
	return "😿"
}

func getMoodText(state State) string {
	if state.Happiness >= 80 {
		return "Мурлычет и трётся об ноги"
	} else if state.Happiness >= 50 {
		return "Играет с клубком ниток"
	} else if state.Happiness >= 30 {
		return "Сидит на подоконнике и смотрит в окно"
	}
	return "Лежит в углу и грустит"
}

func showStatus(state State) {
	fmt.Println(colorize(catASCII, magenta))
	fmt.Println(colorize("  Имя: "+state.Name, blue))
	fmt.Println(colorize("  Возраст: "+strconv.Itoa(state.Age)+" лет", blue))
	fmt.Printf("  %s  %s\n", getEmotion(state), getMoodText(state))
	hungerColor := green
	if state.Hunger > 70 {
		hungerColor = yellow
	}
	fmt.Println(colorize("  Голод: "+strconv.Itoa(state.Hunger)+"/100", hungerColor))
	happyColor := green
	if state.Happiness < 50 {
		happyColor = red
	}
	fmt.Println(colorize("  Счастье: "+strconv.Itoa(state.Happiness)+"/100", happyColor))
	energyColor := green
	if state.Energy < 50 {
		energyColor = yellow
	}
	fmt.Println(colorize("  Энергия: "+strconv.Itoa(state.Energy)+"/100", energyColor))
	healthColor := green
	if state.Health < 30 {
		healthColor = red
	}
	fmt.Println(colorize("  Здоровье: "+strconv.Itoa(state.Health)+"/100", healthColor))
}

func randomEvent(msg string) {
	events := []string{"Мяу!", "Котик трётся о ноги.", "Принёс игрушку.",
		"Пытается поймать муху.", "Свернулся клубком."}
	if msg == "" {
		msg = events[rand.Intn(len(events))]
	}
	fmt.Println(colorize("  ✨ "+msg, yellow))
}

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println(colorize("Usage: virtual_cat <status|feed|play|sleep|heal|rename|auto> [name]", yellow))
		os.Exit(1)
	}
	action := os.Args[1]
	state := loadState()

	switch action {
	case "status":
		showStatus(state)
	case "feed":
		state.Hunger = clamp(state.Hunger-30, 0, 100)
		state.Happiness = clamp(state.Happiness+10, 0, 100)
		state.Health = clamp(state.Health+5, 0, 100)
		fmt.Println(colorize("🐟  Ням-ням! Котик поел.", green))
		randomEvent("Мурлычет от удовольствия")
		saveState(state)
	case "play":
		if state.Energy < 20 {
			fmt.Println(colorize("😿  Котик слишком устал для игр.", red))
			return
		}
		state.Happiness = clamp(state.Happiness+25, 0, 100)
		state.Energy = clamp(state.Energy-20, 0, 100)
		state.Hunger = clamp(state.Hunger+10, 0, 100)
		fmt.Println(colorize("🧶  Игра с клубком! Котик доволен.", green))
		randomEvent("Прыгает за лазерной указкой")
		saveState(state)
	case "sleep":
		state.Energy = clamp(state.Energy+40, 0, 100)
		state.Hunger = clamp(state.Hunger+10, 0, 100)
		fmt.Println(colorize("😴  Котик уснул. Сладких снов!", blue))
		randomEvent("Во сне дёргает лапками")
		saveState(state)
	case "heal":
		if state.Energy < 20 {
			fmt.Println(colorize("😿  Нет сил лечить котика.", red))
			return
		}
		state.Health = clamp(state.Health+30, 0, 100)
		state.Energy = clamp(state.Energy-20, 0, 100)
		fmt.Println(colorize("💊  Котик вылечен! Он благодарен.", green))
		saveState(state)
	case "rename":
		if len(os.Args) < 3 {
			fmt.Println(colorize("Укажите имя: rename <имя>", red))
			return
		}
		state.Name = os.Args[2]
		fmt.Println(colorize("🐱  Котика теперь зовут "+state.Name+"!", blue))
		saveState(state)
	case "auto":
		fmt.Println(colorize("🤖  Автоматический режим включён.", magenta))
		for {
			state.Hunger = clamp(state.Hunger+rand.Intn(21)-5, 0, 100)
			state.Happiness = clamp(state.Happiness+rand.Intn(16)-5, 0, 100)
			state.Energy = clamp(state.Energy+rand.Intn(16)-5, 0, 100)
			state.Health = clamp(state.Health+rand.Intn(9)-3, 0, 100)
			saveState(state)
			clearScreen()
			showStatus(state)
			fmt.Println(colorize("\nНажмите Ctrl+C для выхода", yellow))
			time.Sleep(10 * time.Second)
		}
	default:
		fmt.Println(colorize("Неизвестное действие.", red))
	}
}
