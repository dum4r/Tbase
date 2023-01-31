package core

import (
	"embed"
	"tbase/core/util"
	"tbase/core/w32"
)

const (
	fileConfig = "config.json"
)

var (
	isRun     bool
	isHung    bool
	processes queueProcesses

	opt opts // user options of the window
	win w32.Win32
)

func Start(assets *embed.FS) error {
	util.SetAssets(assets)

	// Create and Save Options for the window
	defer opt.SaveConfigs()
	if err := opt.loadConfigs(); err != nil {
		return err
	}

	// Create window with user32.dll Library
	defer win.Destroy()
	if err := win.Create(&isRun, opt.style, opt.Screen.Rect()); err != nil {
		return err
	}
	// WASAPI api para el sonido del la ventana???????
	// Mmdevapi.dll: Este archivo DLL proporciona funciones de audio para Windows.
	// Audioses.dll: Este archivo DLL contiene funciones que te permiten controlar la sesión de audio de Windows.
	// Mmcsd.dll: Este archivo DLL proporciona funciones para el manejo de dispositivos de audio en Windows.
	// Mmres.dll: Este archivo DLL contiene funciones que te permiten gestionar los recursos de audio de Windows.
	// OOOOOOO tocara con XAudio2
	// El problema es que no existe documentacion directa si no a => audioclient.h

	// Core Loop
	isRun = true
	for isRun {
		// If the window hangs, stop all processes with tickers
		if newState := win.IsHung(); isHung != newState {
			_SetStateProcesses(newState)
		}

		win.Manager()
	}
	win.Wait()
	return nil
}
