package modules

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"context"
	"fmt"
	"github.com/jonas747/dca"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func TextToSpeech(text string, outputFile string) {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},

		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "tr-TR",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},

		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	filename := outputFile
	err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func ConvertDCA(readFile string, outputFile string) {
	var encodeOpts = &dca.EncodeOptions{
		Volume:           256,
		Channels:         2,
		FrameRate:        48000,
		FrameDuration:    20,
		Bitrate:          64,
		Application:      dca.AudioApplicationAudio,
		CompressionLevel: 10,
		PacketLoss:       1,
		BufferedFrames:   100,
		VBR:              true,
		StartTime:        0,
		RawOutput:        true,
	}
	encodeSession, _ := dca.EncodeFile(readFile, encodeOpts)
	defer encodeSession.Cleanup()
	output, err := os.Create(outputFile)

	if err != nil {
		fmt.Println(err)
	}

	io.Copy(output, encodeSession)
}
