package modules

import (
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
	"time"
)

func PlaySound(vc *discordgo.VoiceConnection) (err error) {
	time.Sleep(100 * time.Millisecond)
	vc.Speaking(true)

	for _, buff := range Buffer {
		vc.OpusSend <- buff
	}

	vc.Speaking(false)
	time.Sleep(100 * time.Millisecond)

	return
}

func loadSound(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		Buffer = append(Buffer, InBuf)
	}
}
