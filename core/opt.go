package core

import (
	"encoding/json"
	"os"
	"tbase/core/w32"
)

type opts struct {
	Lang       string
	Fullscreen bool
	Screen     w32.Screen // Position and Dimesion
	TPS        int
	style      int
}

// this function modifies the styles to make the window normal Window

func (o *opts) DefaultStyle() {
	o.Fullscreen = false
	o.Screen.Default() // Set Default position and Dimesion
	o.DefaultScreenStyle()
}
func (o *opts) DefaultScreenStyle() { o.style = w32.WS_MY_OVERLAPPED | w32.WS_VISIBLE }

// this function modifies the styles to make the window fullscreen
func (o *opts) FullScreemStyle() {
	o.Fullscreen = true
	o.Screen.Fullscreen()                         // Set 0,0 and take the whole screen
	o.style = w32.WS_POPUPWINDOW | w32.WS_VISIBLE // Style = POPUPWINDOW
}

// load Config if file-config not exist, create the file with data predefined
func (o *opts) loadConfigs() error {
	if _, err := os.Stat(fileConfig); os.IsNotExist(err) {
		if _, err := os.Create(fileConfig); err != nil {
			return err
		}

		o.TPS = 60
		o.Lang = "es"
		// o.FullScreemStyle()
		o.DefaultStyle()

		file, _ := json.MarshalIndent(o, "", " ")
		if err := os.WriteFile(fileConfig, file, 0644); err != nil {
			return err
		}
		return nil
	} else {
		dat, err := os.ReadFile(fileConfig)
		if err != nil {
			return err
		}
		if len(dat) < 1 {
			panic("Error novel: no data in file loading config")
		}
		if err := json.Unmarshal([]byte(dat), &o); err != nil {
			print("Error novel: corrupted file")
			panic(err)
		}
	}
	if o.Fullscreen {
		o.FullScreemStyle()
	} else {
		o.DefaultScreenStyle()
	}
	return nil
}

func (o *opts) SaveConfigs() error {
	file, _ := json.MarshalIndent(o, "", " ")
	if err := os.WriteFile(fileConfig, file, 0644); err == nil {
		return err
	}
	return nil
}
