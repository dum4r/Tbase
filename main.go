package main

import (
	"embed"
	_ "image/png"
	"math/rand"
	"tbase/core"
	"time"
)

var (
	//go:embed assets/*
	assets embed.FS
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	if err := core.Start(&assets); err != nil {
		panic(err)
	}
}
