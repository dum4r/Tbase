package w32

func (w *Win32) _Buffer(widthBuff, heightBuff int, ix, iy int) []byte {
	widthScreen := w.screen.Width()
	if widthScreen < widthBuff+ix {
		panic("Error: Buffer Width out limit!")
	}

	if w.screen.Height() < heightBuff+iy {
		panic("Error: Buffer Height out limit!")
	}

	newBuff := make([]byte, widthBuff*heightBuff*4)
	for y := 0; y < heightBuff; y++ {
		for x := 0; x < widthBuff; x++ {
			indexSR := (y*widthBuff*4 + x*4)
			indexRR := ((y + iy) * widthScreen * 4) + ((x + ix) * 4)
			newBuff[indexSR] = w.buf[indexRR]     // b
			newBuff[indexSR+1] = w.buf[indexRR+1] // g
			newBuff[indexSR+2] = w.buf[indexRR+2] // r
		}
	}
	return newBuff
}
