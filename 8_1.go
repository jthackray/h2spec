package h2spec

import (
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func TestHTTPRequestResponseExchange(ctx *Context) {
	if !ctx.IsTarget("8.1") {
		return
	}

	PrintHeader("8.1. HTTP Request/Response Exchange", 0)
	TestHTTPHeaderFields(ctx)
	PrintFooter()
}

func TestHTTPHeaderFields(ctx *Context) {
	PrintHeader("8.1.2. HTTP Header Fields", 1)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains the header field name in uppercase letters"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("X-TEST", "test"),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 1)
	}(ctx)

	TestPseudoHeaderFields(ctx)
	TestConnectionSpecificHeaderFields(ctx)
	TestRequestPseudoHeaderFields(ctx)
	TestMalformedRequestsAndResponses(ctx)
}

func TestPseudoHeaderFields(ctx *Context) {
	PrintHeader("8.1.2.1. Pseudo-Header Fields", 2)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains the pseudo-header field defined for response"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair(":status", "200"),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains the invalid pseudo-header field"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair(":test", "test"),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains a pseudo-header field that appears in a header block after a regular header field"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair("x-test", "test"),
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)
}

func TestConnectionSpecificHeaderFields(ctx *Context) {
	PrintHeader("8.1.2.2. Connection-Specific Header Fields", 2)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains the connection-specific header field"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("connection", "keep-alive"),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains the TE header field that contain any value other than \"trailers\""
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("trailers", "test"),
			pair("te", "trailers, deflate"),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)
}

func TestRequestPseudoHeaderFields(ctx *Context) {
	PrintHeader("8.1.2.3. Request Pseudo-Header Fields", 2)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that is omitted mandatory pseudo-header fields"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":authority", ctx.Authority()),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)
}

func TestMalformedRequestsAndResponses(ctx *Context) {
	PrintHeader("8.1.2.6. Malformed Requests and Responses", 2)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains invalid \"content-length\" header field"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair(":method", "POST"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("content-length", "1"),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = false
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)
		http2Conn.fr.WriteData(1, true, []byte("test"))

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)
}
