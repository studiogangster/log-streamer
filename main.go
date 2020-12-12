package main

import (
	"flag"
	"fmt"
	"log"
	"log_reader/readlogs"
	"log_reader/sanitize"
	"strconv"

	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8082", "TCP address to listen to")
	compress = flag.Bool("compress", true, "Whether to enable transparent response compression")
)

func main() {
	flag.Parse()

	h := func(ctx *fasthttp.RequestCtx) {

		defer func() {
			if r := recover(); r != nil {
				log.Println("Somethig webt wron", r)
				ctx.Error("Error reading logs", 403)

			}
		}()
		requestHandler(ctx)
	}

	// h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {

	fileID := ctx.Request.Header.Peek("X-Id")

	fileName := sanitize.BaseName(string(fileID)) + ".log"

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
