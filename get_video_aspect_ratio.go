package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	buffer := bytes.Buffer{}
	cmd.Stdout = &buffer
	if err := cmd.Run(); err != nil {
		return fmt.Sprintln("Unable to run command"), err
	}

	type parameters struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	var params parameters

	if err := json.Unmarshal(buffer.Bytes(), &params); err != nil {
		return fmt.Sprintln("Unable to unmarshal buffer"), err
	}

	width := float64(params.Streams[0].Width)
	height := float64(params.Streams[0].Height)

	ratio := width / height
	targetLandscape := 16.0 / 9.0
	targetPortrait := 9.0 / 16.0
	tolerance := 0.2

	switch {
	case math.Abs(ratio-targetLandscape) < tolerance:
		return "16:9", nil
	case math.Abs(ratio-targetPortrait) < tolerance:
		return "9:16", nil
	default:
		return "other", nil
	}
}
