package main

import (
	"fmt"
	"log"
	"log_reader/readlogs"
	"log_reader/sanitize"
	"os"
	"path/filepath"
	"strconv"

	"github.com/valyala/fasthttp"
)

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

var (
	addr       = getEnv("addr", ":8082")
	logdir     = getEnv("log_dir", "./test_log_files")
	filesuffix = getEnv("suffix", "log")
	fileprefix = getEnv("prefix", "filewatchdog")
)

func main() {

	// content, err := readlogs.Read("secoflex-localized-logs/watchdog.cortx.nucleus/log_", 5, 20000, false)
	content, err := readlogs.Read("secoflex-localized-logs/watchdog.6696ddea59f1fb1cd57d885a/log_", 5, 20000, false)

	log.Println(string(content), err)
	return

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
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")

}

func requestHandler(ctx *fasthttp.RequestCtx) {

	fileID := ctx.Request.Header.Peek("X-Id")

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

	buffer, startByte, endByte, fileSize := readlogs.StreamFile(fileName, seekTo, maxLineCount)
	fmt.Fprintf(ctx, string(buffer))
	ctx.SetContentType("text/plain; charset=utf8")
	ctx.Response.Header.Set("X-File-Range", fmt.Sprintf("%d-%d", startByte, endByte))
	ctx.Response.Header.Set("X-File-Size", fmt.Sprintf("%d", fileSize))

}
