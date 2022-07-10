package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	_ = os.Mkdir(fmt.Sprintf("%s", "./media"), os.ModePerm)
	Cmd()
}

func GinMode() string {
	gin_mode := getEnv("GIN_MODE", "debug")
	return fmt.Sprintf("%s", gin_mode)
}

func Url() string {
	url := getEnv("URL", "http://streamyes.alsolnet.com/buturama/live/playlist.m3u8")
	return fmt.Sprintf("%s", url)
}

func Name() string {
	name := getEnv("NAME", "channel")
	return fmt.Sprintf("%s", name)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func Cmd() {

	path := fmt.Sprintf("./media/%s", Name())

	args := []string{"-f", "hls", "-i", Url()}
	args = append(args, "-map", "0:v:0", "-map", "0:a:0", "-map", "0:v:0", "-map", "0:a:0", "-map", "0:v:0", "-map", "0:a:0")
	args = append(args, "-c:v", "libx264", "-preset:v", "ultrafast", "-c:a", "copy", "-x264opts", "keyint=48:min-keyint=48:no-scenecut", "-sc_threshold", "0")
	args = append(args, "-b:v:0", "500k", "-maxrate:v:0", "600k", "-bufsize:v:0", "500k")
	args = append(args, "-b:v:1", "1500k", "-maxrate:v:1", "1800k", "-bufsize:v:1", "1500k")
	args = append(args, "-b:v:2", "3000k", "-maxrate:v:2", "3500k", "-bufsize:v:2", "3000k")
	args = append(args, "-var_stream_map", "v:0,a:0,name:480p v:1,a:1,name:720p v:2,a:2,name:1080p")
	args = append(args, "-f", "hls", "-hls_list_size", "6", "-threads", "0", "-hls_time", "3", "-hls_flags", "delete_segments")
	args = append(args, "-master_pl_name", "playlist.m3u8", "-y", fmt.Sprintf("%s/%s", path, "stream-%v.m3u8"))

	for {
		_ = os.Mkdir(fmt.Sprintf("%s", path), os.ModePerm)
		cmd := exec.Command("ffmpeg", args...)
		if err := cmd.Start(); err != nil {
			log.Printf("Command start with error: %v", err)
		}
		err := cmd.Wait()
		log.Printf("Command finished with error: %v", err)
		cmdKill := exec.Command("pkill", "-f", "ffmpeg")
		_ = cmdKill.Start()
		_ = cmdKill.Wait()
		_ = os.RemoveAll(fmt.Sprintf("%s", path))
		time.Sleep(2 * time.Second)
	}
}
