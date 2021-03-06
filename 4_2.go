package h2spec

import (
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func TestFrameSize(ctx *Context) {
	if !ctx.IsTarget("4.2") {
		return
	}

	PrintHeader("4.2. Frame Size", 0)
	msg := "The endpoint MUST send a FRAME_SIZE_ERROR error."

	func(ctx *Context) {
		desc := "Sends large size frame that exceeds the SETTINGS_MAX_FRAME_SIZE"
		result := false

		http2Conn := CreateHttp2Conn(ctx, false)
		defer http2Conn.conn.Close()

		http2Conn.fr.WriteSettings()

		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = false
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)
		http2Conn.fr.WriteData(1, true, []byte(GetDummyData(16385)))

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeFrameSize {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
