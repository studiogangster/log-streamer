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
)

func main() {

	h := func(ctx *fasthttp.RequestCtx) {

		defer func() {
			if r := recover(); r != nil {
				log.Println("Somethig webt wron", r)
				ctx.Error("Error reading logs", 403)

			}
		}()
		requestHandler(ctx)
	}

	h = fasthttp.CompressHandler(h)

	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {

	fileID := ctx.Request.Header.Peek("X-Id")

	fileName := sanitize.BaseName(string(fileID)) + "." + filesuffix

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
