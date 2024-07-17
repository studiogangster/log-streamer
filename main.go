package main

import (
	"fmt"
	"log"
	"log_reader/readlogs"
	"log_reader/sanitize"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

var logdir string
var (
	addr       = getEnv("addr", ":8082")
	filesuffix = getEnv("suffix", "log")
	fileprefix = getEnv("prefix", "filewatchdog")
)

func readFromMinio() {

	// content, err := readlogs.Read("secoflex-localized-logs/watchdog.cortx.nucleus/log_", 5, 20000, false)
	content, fileSize, err := readlogs.Read(logdir+"/watchdog.6696fa1059f1fb1cd57d885b/log_", 0, 400, false)

	log.Println(string(content), fileSize, err)
	return
}

func main() {
	godotenv.Load()
	logdir = getEnv("LOG_DIR", "")

	// readFromMinio()

	h := func(ctx *fasthttp.RequestCtx) {

		setCorsHeaders(ctx)

		defer func() {
			if r := recover(); r != nil {
				log.Println("Somethig webt wron", r)
				setCorsHeaders(ctx)

				// ctx.Request.Body("Error reading logs", 200)

			}
		}()

		if ctx.IsGet() {
			requestHandler(ctx)

		} else {

		}

	}

	h = fasthttp.CompressHandler(h)

	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func setCorsHeaders(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
	ctx.Response.Header.Set("Access-Control-Expose-Headers", "X-File-Range, X-File-Size")

	ctx.Response.Header.Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")

}

func requestHandler(ctx *fasthttp.RequestCtx) {

	fileID := string(ctx.Request.Header.Peek("X-Id"))

	fileName := fileprefix + "." + sanitize.BaseName(string(fileID)) + "." + filesuffix

	fileName = filepath.Join(logdir, fileName)
	log.Println(fileName)

	seekTo, err := strconv.ParseInt(string(ctx.Request.Header.Peek("X-Seek")), 10, 0)

	if err != nil {
		seekTo = 0
	}

	maxLineCount, err := strconv.ParseInt(string(ctx.Request.Header.Peek("X-Lines")), 10, 0)

	if err != nil {
		maxLineCount = 100
	}

	// buffer, startByte, endByte, fileSize := readlogs.StreamFile(fileName, seekTo, maxLineCount)

	buffer, fileSize, err := readlogs.Read(logdir+"/watchdog."+fileID+"/log_", seekTo, maxLineCount, false)
	log.Println(string(buffer), fileSize, err)

	fmt.Fprintf(ctx, string(buffer))
	ctx.SetContentType("text/plain; charset=utf8")
	ctx.Response.Header.Set("X-File-Range", fmt.Sprintf("%d-%d", seekTo, seekTo+int64(len(buffer))))
	ctx.Response.Header.Set("X-File-Size", fmt.Sprintf("%d", fileSize))

}
