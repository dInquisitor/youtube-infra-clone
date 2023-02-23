package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/segmentio/kafka-go"
)

const VIDEO_PROCESSING_TOPIC = "process-video"
const VIDEO_PROCESSING_PARTITION = 0
const VIDEO_PROCESSING_GROUP = "process-video-group"

const VIDEO_DIRECTORY = "../data/video-store"

func main() {
	log.Println("Running video processing worker")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     VIDEO_PROCESSING_TOPIC,
		Partition: VIDEO_PROCESSING_PARTITION,
		GroupID:   VIDEO_PROCESSING_GROUP,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	fmt.Println("connected to kafka")

	defer func() {
		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}

		go processVideo(string(m.Value))
	}
}

func processVideo(videoID string) {
	resolutions := map[string][2]string{
		// height, width and bit rate
		"480": {"640", "250k"},
		// "720":  {"1280", "750k"},
		// "1080": {"1920", "1000k"},
		// "2160": {"3840", "1500k"},
	}

	// originalFilePath, err := getAbsPath(videoID, "original")

	// if err != nil {
	// 	return
	// }

	// var wg sync.WaitGroup
	// // process all resolutions and audio stream in parallel
	// wg.Add(len(resolutions) + 1)

	// go makeAudio(originalFilePath, videoID)

	// for height, details := range resolutions {
	// 	go makeVideoWithResolution(originalFilePath, videoID, details[0], height, details[1])
	// }

	// wg.Wait()

	// generate manifest after generating audio and video
	makeManifest(resolutions, videoID)
}

func makeManifest(resolutions map[string][2]string, videoID string) {
	var arguments []string
	var adaptSet []string

	for height := range resolutions {
		videoFilePath, err := getAbsPath(videoID, fmt.Sprintf("%s.webm", height))

		if err != nil {
			return
		}

		arguments = append(arguments, "-f", "webm_dash_manifest", "-i", videoFilePath)
	}

	audioFilePath, err := getAbsPath(videoID, "audio.webm")

	if err != nil {
		return
	}

	arguments = append(arguments, "-f", "webm_dash_manifest", "-i", audioFilePath)
	arguments = append(arguments, "-c", "copy")

	lastI := 0
	for i := 0; i < len(resolutions); i++ {
		arguments = append(arguments, "-map", fmt.Sprintf("%d:0", i))
		adaptSet = append(adaptSet, fmt.Sprintf("%d", i))
		lastI = i
	}

	arguments = append(arguments, "-map", fmt.Sprintf("%d:0", lastI+1))
	arguments = append(arguments, "-y", "-copy_unknown")
	arguments = append(arguments, "-f", "webm_dash_manifest", "-adaptation_sets", fmt.Sprintf("id=0,streams=%s id=1,streams=%d", strings.Join(adaptSet, ","), lastI+1))

	manifestFilePath, err := getAbsPath(videoID, "manifest.mpd")

	if err != nil {
		return
	}

	arguments = append(arguments, manifestFilePath)

	log.Println(arguments)
	log.Println(strings.Join(arguments, ","))

	runFFmpeg(arguments)
}

func makeAudio(originalFilePath, videoID string) {
	audioFilePath, err := getAbsPath(videoID, "audio.webm")

	if err != nil {
		return
	}

	// ffmpeg -i original -vn -acodec libvorbis -ab 128k -dash 1 audio.webm
	runFFmpeg([]string{"-i", originalFilePath, "-vn", "-acodec", "libvorbis", "-ab", "128k", "-dash", "1", "-y", audioFilePath})
}

func makeVideoWithResolution(originalFilePath, videoID, width, height, bitRate string) {
	destinationFilePath, err := getAbsPath(videoID, fmt.Sprintf("%s.webm", height))
	if err != nil {
		return
	}

	runFFmpeg([]string{"-i", originalFilePath, "-c:v", "libvpx-vp9", "-keyint_min", "150", "-g", "150", "-tile-columns", "4", "-frame-parallel", "1", "-f", "webm", "-dash", "1", "-an", "-vf", fmt.Sprintf("scale=%s:%s", width, height), "-b:v", bitRate, "-dash", "1", "-y", destinationFilePath})
}

func runFFmpeg(arguments []string) {
	arguments = append(arguments, "-loglevel", "repeat+level+verbose")
	cmd := exec.Command("ffmpeg", arguments...)

	fmt.Printf("cmd: %v\n", cmd)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Println(err)
		log.Println(stderr.String())
	}

	fmt.Printf("Done running cmd: %v\n", cmd)
}

func getAbsPath(videoID string, filename string) (string, error) {
	absFilePath, err := filepath.Abs(fmt.Sprintf("%s/%s/%s", VIDEO_DIRECTORY, videoID, filename))

	if err != nil {
		log.Println("Could not expand file path", err)
		return "", err
	}

	return absFilePath, nil
}
