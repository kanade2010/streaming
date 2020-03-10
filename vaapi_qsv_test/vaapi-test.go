package main

import(
	"github.com/ailumiyana/goav-incr/goav/avcodec"
	"github.com/ailumiyana/goav-incr/goav/avutil"
	"github.com/ailumiyana/goav-incr/goav/avfilter"

	"github.com/ailumiyana/tools/latency"
	log "github.com/astaxie/beego/logs"

	"fmt"
	//"time"
	"io/ioutil"
	"os"
	"strconv"
	"unsafe"
)


var input_h264  string = "1080p.h264"
var output_h264 string = "out.h264"

type VaapiHWDeviceCtx struct {
	hw_device_ctx *avutil.BufferRef
}

// input : "/dev/dri/card0" or "/dev/dri/renderD128"
func Create(device string) *VaapiHWDeviceCtx {
	var hw_device_ctx *avutil.BufferRef
	err := avutil.AVHwdeviceCtxCreate(&hw_device_ctx, avutil.AV_HWDEVICE_TYPE_VAAPI,
		device, nil, 0)
	if err < 0 {
		fmt.Println("AVHwdeviceCtxCreate err : ", avutil.ErrorFromCode(err))
		return nil
	}

	return &VaapiHWDeviceCtx{
		hw_device_ctx,
	}
}

func (v *VaapiHWDeviceCtx)Context() *avutil.BufferRef{
	return v.hw_device_ctx
}

//warning need unref by avutil.AVBufferUnref()
func (v *VaapiHWDeviceCtx)GetAnAllocHwframeCtxRef(format, sw_format avutil.PixelFormat, w, h int) *avutil.BufferRef{
	hw_frames_ref := avutil.AVHwframeCtxAlloc(v.hw_device_ctx)
	if hw_frames_ref == nil {
		fmt.Println("AVHwframeCtxAlloc err")
		return nil
	}

	frames_ctx := hw_frames_ref.HWFramesContext()
	//frames_ctx.SetHWFramesContextPrarms(avutil.PixelFormat(avcodec.AV_PIX_FMT_VAAPI),
	//		avutil.PixelFormat(avcodec.AV_PIX_FMT_YUV420P),	
	//		1280, 720)
	frames_ctx.SetHWFramesContextPrarms(format, sw_format, w, h)

	err := avutil.AVHwFrameCtxInit(hw_frames_ref)
	if err < 0 {
		fmt.Println("AVHwFrameCtxInit err : ", avutil.ErrorFromCode(err))
		avutil.AVBufferUnref(&hw_frames_ref)
		return nil
	}

	return hw_frames_ref
}

//var hw_device_ctx *avutil.BufferRef// = avutil.AVHwdeviceCtxAlloc(avutil.AV_HWDEVICE_TYPE_VAAPI)

func SetHwframesContext(ctxt *avcodec.Context) {
	/*var hw_device_ctx *avutil.BufferRef

	err := avutil.AVHwdeviceCtxCreate(&hw_device_ctx, avutil.AV_HWDEVICE_TYPE_VAAPI,
		"/dev/dri/card0", nil, 0)
	if err < 0 {
		fmt.Println("AVHwdeviceCtxCreate err : ", avutil.ErrorFromCode(err))
		panic("AVHwdeviceCtxCreate faild")
	}*/

	vaapi_device  := Create("/dev/dri/card0")
	//hw_device_ctx := vaapi_device.Context()

	/*hw_frames_ref := avutil.AVHwframeCtxAlloc(hw_device_ctx)
	if hw_frames_ref == nil {
		fmt.Println("AVHwframeCtxAlloc err")
		panic("AVHwframeCtxAlloc faild")
	}

	frames_ctx := hw_frames_ref.HWFramesContext()
	frames_ctx.SetHWFramesContextPrarms(avutil.PixelFormat(avcodec.AV_PIX_FMT_VAAPI),
			avutil.PixelFormat(avcodec.AV_PIX_FMT_YUV420P),	
			1280, 720)

	err := avutil.AVHwFrameCtxInit(hw_frames_ref)
	if err < 0 {
		fmt.Println("AVHwFrameCtxInit err : ", avutil.ErrorFromCode(err))
		avutil.AVBufferUnref(&hw_frames_ref)
		panic("AVHwFrameCtxInit faild")
	}*/

	hw_frames_ref := vaapi_device.GetAnAllocHwframeCtxRef(
							avutil.PixelFormat(avcodec.AV_PIX_FMT_VAAPI),
							avutil.PixelFormat(avcodec.AV_PIX_FMT_YUV420P),	1280, 720)

	//enc.Context().HWFramesCtx() = avutil.AVBufferRef(hw_frames_ref)
	ctxt.SetHWFramesCtx(avutil.AVBufferRef(hw_frames_ref))
	if ctxt.HWFramesCtx() == nil {
		fmt.Println("SetHWFrameCtx err ")
		avutil.AVBufferUnref(&hw_frames_ref)
		panic("SetHWFrameCtx faild")
	}

	//fmt.Println(frames_ctx)
	//fmt.Println(hw_frames_ref)
	//fmt.Println(ctxt.HWFramesCtx())

	avutil.AVBufferUnref(&hw_frames_ref)
}


func decode_encode() {
	//decoder
	vaapi_device  := Create("/dev/dri/renderD128")

	pkt 		   := avcodec.AvPacketAlloc()
	if pkt == nil {
		log.Critical("AvPacketAlloc failed.")
		return 
	}

	codec 		   := avcodec.AvcodecFindDecoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
	if codec == nil {
		log.Critical("AvcodecFindDecoder failed.")
	}

	context 	   := codec.AvcodecAllocContext3()
	if context == nil {
		log.Critical("AvcodecAllocContext3 failed.")
		return 
	}

	parserContext  := avcodec.AvParserInit(int(avcodec.CodecId(avcodec.AV_CODEC_ID_H264)))
	if parserContext == nil {
		log.Critical("AvParserInit failed.")
		return 
	}

	frame   	   := avutil.AvFrameAlloc()
	if frame == nil {
		log.Critical("AvFrameAlloc failed.")
		return 
	}

	context.SetHWDeviceCtx(avutil.AVBufferRef(vaapi_device.Context()))
	err := context.AvcodecOpen2(codec, nil)
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return 
	}

	// encoder
	codec_enc 		   := avcodec.AvcodecFindEncoderByName("h264_vaapi")
	if codec_enc == nil {
		log.Critical("AvcodecFindEncoderByName failed.")
	}

	pkt_enc 		   := avcodec.AvPacketAlloc()
	if pkt_enc == nil {
		log.Critical("AvPacketAlloc failed.")
	}

	context_enc 	   := codec_enc.AvcodecAllocContext3()
	if context_enc == nil {
		log.Critical("AvcodecAllocContext3 failed.")
	}

	context_enc.SetVideoEncodeParams(2000000, 1920, 1080,
		avcodec.AV_PIX_FMT_VAAPI,
		false, 10)
	context_enc.SetTimebase(1, 30)

	hw_frames_ref := vaapi_device.GetAnAllocHwframeCtxRef(
							avutil.PixelFormat(avcodec.AV_PIX_FMT_VAAPI),
							avutil.PixelFormat(avcodec.AV_PIX_FMT_NV12), 1920, 1080)


	context_enc.SetHWFramesCtx(avutil.AVBufferRef(hw_frames_ref))

	err = context_enc.AvcodecOpen2(codec_enc, nil)
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return
	}

	decLatency := latency.New("vaapi", "decode")
	encLatency := latency.New("vaapi", "encode")
	/*scaleLatency := latency.New("vaapi", "scale")
	overlayLatency := latency.New("vaapi", "overlay")*/

	data, erro := ioutil.ReadFile(input_h264)
    if erro != nil {
        log.Debug("File reading error", erro)
        return
	}
	log.Debug("Open Success.")
	l := len(data)
    log.Debug("size of file:", l)
	
	b := make([]byte, 4096 + 64)
	
	file, erro := os.Create(output_h264)
	if erro != nil {
		log.Critical("Error Reading")
	}
	defer file.Close()

	//var pts int64 = 0
	//var opts int64 = 0

	//var frames [25]*avutil.Frame
	started := false
	sum := 0
	for sum < l {
		remain := 4096
		for remain > 0 {
			copy(b, data[sum:sum + 4096])
			if !started {
				decLatency.Start()
				started = true
			}
			n := context.AvParserParse2(parserContext, pkt, b, 
				remain, avcodec.AV_NOPTS_VALUE, avcodec.AV_NOPTS_VALUE, 0)
			log.Debug("parser ", n, "bytes")
			
			sum     = sum + n
			remain  = remain - n;

			//log.Trace("--------", dec.Packet().GetPacketSize())
			if pkt.GetPacketSize() > 0 {
				//fmt.Println(decLatency.End())

				ret := context.AvcodecSendPacket(pkt)
				if ret < 0 {
					log.Error("AvcodecSendPacket err ", avutil.ErrorFromCode(ret))
					continue 
				}
			
				ret = context.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(frame)))
				if ret < 0 {
					log.Error("AvcodecReceiveFrame err ", avutil.ErrorFromCode(ret))
					continue 
				}

				if ret == 0 {
					started = false
					decLatency.End()
					fmt.Println(decLatency.End())

					log.Debug("dec-frame:",frame)
					

					encLatency.Start()
					ret = context_enc.AvcodecSendFrame((*avcodec.Frame)(unsafe.Pointer(frame)))
					if ret < 0 {
						log.Trace("AvcodecSendFrame err ", avutil.ErrorFromCode(ret))
						continue
					}
				
					ret = context_enc.AvcodecReceivePacket(pkt_enc)
					if ret < 0 {
						log.Trace("AvcodecReceivePacket err ", avutil.ErrorFromCode(ret))
						continue
					}
					if ret == 0 {
						data0 := pkt_enc.Data()
						buf := make([]byte, pkt_enc.GetPacketSize())
						start := uintptr(unsafe.Pointer(data0))
						for i := 0; i < pkt_enc.GetPacketSize(); i++ {
							elem := *(*uint8)(unsafe.Pointer(start + uintptr(i)))
							buf[i] = elem
						}
						
						file.Write(buf)
						//encLatency.End()
						fmt.Println(encLatency.End())
					}

					avutil.AvFrameUnref(frame)
					//avutil.AvFrameUnref(oframe)

				}

			}
		}
	}	
}

func decode_overlay_encode() {
	//decoder
	vaapi_device  := Create("/dev/dri/renderD128")

	pkt 		   := avcodec.AvPacketAlloc()
	if pkt == nil {
		log.Critical("AvPacketAlloc failed.")
		return 
	}

	codec 		   := avcodec.AvcodecFindDecoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
	if codec == nil {
		log.Critical("AvcodecFindDecoder failed.")
	}

	context 	   := codec.AvcodecAllocContext3()
	if context == nil {
		log.Critical("AvcodecAllocContext3 failed.")
		return 
	}

	parserContext  := avcodec.AvParserInit(int(avcodec.CodecId(avcodec.AV_CODEC_ID_H264)))
	if parserContext == nil {
		log.Critical("AvParserInit failed.")
		return 
	}

	frame   	   := avutil.AvFrameAlloc()
	if frame == nil {
		log.Critical("AvFrameAlloc failed.")
		return 
	}

	context.SetHWDeviceCtx(avutil.AVBufferRef(vaapi_device.Context()))
	err := context.AvcodecOpen2(codec, nil)
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return 
	}

	// encoder
	codec_enc 		   := avcodec.AvcodecFindEncoderByName("h264_vaapi")
	if codec_enc == nil {
		log.Critical("AvcodecFindEncoderByName failed.")
	}

	pkt_enc 		   := avcodec.AvPacketAlloc()
	if pkt_enc == nil {
		log.Critical("AvPacketAlloc failed.")
	}

	context_enc 	   := codec_enc.AvcodecAllocContext3()
	if context_enc == nil {
		log.Critical("AvcodecAllocContext3 failed.")
	}

	context_enc.SetVideoEncodeParams(2000000, 1920, 1080,
		avcodec.AV_PIX_FMT_VAAPI,
		false, 10)
	context_enc.SetTimebase(1, 30)

	hw_frames_ref := vaapi_device.GetAnAllocHwframeCtxRef(
							avutil.PixelFormat(avcodec.AV_PIX_FMT_VAAPI),
							avutil.PixelFormat(avcodec.AV_PIX_FMT_NV12), 1920, 1080)


	/*context_enc.SetHWFramesCtx(avutil.AVBufferRef(outs[0].AvBuffersinkGetHwFramesCtx()))

	err = context_enc.AvcodecOpen2(codec_enc, nil)
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return
	}*/

	decLatency := latency.New("vaapi", "decode")
	encLatency := latency.New("vaapi", "encode")

	// filter
	args := "video_size=1920x1080:pix_fmt=46:time_base=1/30:pixel_aspect=1/1"
	des  := "[in0] fps=30,scale_vaapi=w=960:h=540:format=nv12, hwdownload,format=nv12 [s0];color=black:r=30:size=1920x1080:sar=1/1 [b0];[b0][s0] overlay=x=0:y=0, format=nv12, hwupload"
	
	graph := avfilter.AvfilterGraphAlloc()
	if graph == nil {
		log.Critical("AvfilterGraphAlloc Failed.")
		return
	}
	//graph.SetNbThreads(8)

	inputs  := avfilter.AvfilterInoutAlloc()
	outputs := avfilter.AvfilterInoutAlloc()
	if inputs == nil || outputs == nil {
		log.Critical("AvfilterInoutAlloc Failed.")
		return
	}

	defer avfilter.AvfilterInoutFree(inputs)
	defer avfilter.AvfilterInoutFree(outputs)

	var buffersrc *avfilter.Filter
	var buffersink *avfilter.Filter
	if false {
		buffersrc  = avfilter.AvfilterGetByName("abuffer")
		buffersink = avfilter.AvfilterGetByName("abuffersink")
	} else {
		buffersrc  = avfilter.AvfilterGetByName("buffer")
		buffersink = avfilter.AvfilterGetByName("buffersink")
	}
	
	if buffersink == nil || buffersrc == nil {
		log.Critical("AvfilterGetByName Failed.")
		return
	}

	ret := graph.AvfilterGraphParse2(des, &inputs, &outputs)
	if ret < 0 {
		log.Critical("AvfilterInoutAlloc Failed des : ", avutil.ErrorFromCode(ret))
		return
	}

	if vaapi_device.Context() != nil {
		for _, v := range graph.Filters() {
			v.SetHWDeviceCtx(avutil.AVBufferRef(vaapi_device.Context()))
		}
	}

	var ins    []*avfilter.Context
	var outs   []*avfilter.Context
	var frames []*avutil.Frame
	
	// inputs
	index := 0
	
	for cur := inputs; cur != nil; cur = cur.Next() {
		//log.Debug("index :", index)
		var in *avfilter.Context
		//var args = "video_size=1280x720:pix_fmt=0:time_base=1/30:pixel_aspect=1/1"
		inName := "in" + strconv.Itoa(index)
		ret = avfilter.AvfilterGraphCreateFilter(&in, buffersrc, inName, args, 0, graph)
		if ret < 0 {
			log.Critical("AvfilterGraphCreateFilter Failed des : ", avutil.ErrorFromCode(ret))
			return
		}

		par := avfilter.AvBuffersrcParametersAlloc()
		if par == nil {
			log.Critical("AvBuffersrcParametersAlloc Failed.")
			return
		}
		par.SetHwFramesCtx(avutil.AVBufferRef(hw_frames_ref))
		ret = in.AvBuffersrcParametersSet(par)
		if ret < 0 {
			log.Critical("AvBuffersrcParametersSet Failed.")
			return
		}
		avutil.AvFreep(unsafe.Pointer(&par))

		ins = append(ins, in)
		ret = avfilter.AvfilterLink(ins[index], 0, cur.FilterContext(), cur.PadIdx())
		if ret < 0 {
			log.Critical("AvfilterLink Failed des : ", avutil.ErrorFromCode(ret))
			return
		}
		index++
	}

	// outputs
	index = 0
	for cur := outputs; cur != nil; cur = cur.Next() {
		var out *avfilter.Context
		outName := "out" + strconv.Itoa(index)
		ret = avfilter.AvfilterGraphCreateFilter(&out, buffersink, outName, "", 0, graph)
		if ret < 0 {
			log.Critical("AvfilterGraphCreateFilter Failed des : ", avutil.ErrorFromCode(ret))
			return
		}
	
		outs = append(outs, out)
		ret = avfilter.AvfilterLink(cur.FilterContext(), cur.PadIdx(), outs[index], 0)
		if ret < 0 {
			log.Critical("AvfilterLink Failed des : ", avutil.ErrorFromCode(ret))
			return
		}
		index++

		f := avutil.AvFrameAlloc()
		if f == nil {
			log.Critical("AvFrameAlloc failed.")
			return
		}
		frames = append(frames, f)
	}

	ret = graph.AvfilterGraphConfig(0)
	if ret < 0 {
		log.Critical("AvfilterGraphConfig Failed des : ", avutil.ErrorFromCode(ret))
		return
	}

	overlayLatency := latency.New("vaapi", "overlay")


	context_enc.SetHWFramesCtx(avutil.AVBufferRef(outs[0].AvBuffersinkGetHwFramesCtx()))

	err = context_enc.AvcodecOpen2(codec_enc, nil)
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return
	}

	data, erro := ioutil.ReadFile(input_h264)
    if erro != nil {
        log.Debug("File reading error", erro)
        return
	}
	log.Debug("Open Success.")
	l := len(data)
    log.Debug("size of file:", l)
	
	b := make([]byte, 4096 + 64)
	
	file, erro := os.Create(output_h264)
	if erro != nil {
		log.Critical("Error Reading")
	}
	defer file.Close()

	var pts int64 = 0
	//var opts int64 = 0

	//var frames [25]*avutil.Frame
	started := false
	sum := 0
	for sum < l {
		remain := 4096
		for remain > 0 {
			copy(b, data[sum:sum + 4096])
			if !started {
				decLatency.Start()
				started = true
			}
			n := context.AvParserParse2(parserContext, pkt, b, 
				remain, avcodec.AV_NOPTS_VALUE, avcodec.AV_NOPTS_VALUE, 0)
			log.Debug("parser ", n, "bytes")
			
			sum     = sum + n
			remain  = remain - n;

			//log.Trace("--------", dec.Packet().GetPacketSize())
			if pkt.GetPacketSize() > 0 { // decode
				//fmt.Println(decLatency.End())

				ret := context.AvcodecSendPacket(pkt)
				if ret < 0 {
					log.Error("AvcodecSendPacket err ", avutil.ErrorFromCode(ret))
					continue 
				}
			
				ret = context.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(frame)))
				if ret < 0 {
					log.Error("AvcodecReceiveFrame err ", avutil.ErrorFromCode(ret))
					continue 
				}

				if ret == 0 {
					decLatency.End()
					fmt.Println(decLatency.End())
				} else {
					continue
				}

				overlayLatency.Start()
				pts++
				frame.SetPts(pts)
				
				ret = avfilter.AvBuffersrcAddFrame(ins[0], (*avfilter.Frame)(unsafe.Pointer(frame)))
				if ret < 0 {
					fmt.Println("AvBuffersrcAddFrame error,", avutil.ErrorFromCode(ret))
					continue
				}

				ret = avfilter.AvBufferSinkGetFrame(outs[0], (*avfilter.Frame)(unsafe.Pointer(frames[0])))

				if ret == avutil.AvErrorEOF || ret == avutil.AvErrorEAGAIN {
					log.Error(avutil.ErrorFromCode(ret))
					continue
				}
			
				if ret < 0 {
					log.Error("AvBufferSinkGetFrame Failed des : ", ret, avutil.ErrorFromCode(ret))
					continue
				}
			
				frame = frames[0]

				fmt.Println(overlayLatency.End())

				if ret == 0 { // encode
					started = false

					//log.Debug("dec-frame:",frame)
					
					encLatency.Start()
					ret = context_enc.AvcodecSendFrame((*avcodec.Frame)(unsafe.Pointer(frame)))
					if ret < 0 {
						log.Trace("AvcodecSendFrame err ", avutil.ErrorFromCode(ret))
						continue
					}
				
					ret = context_enc.AvcodecReceivePacket(pkt_enc)
					if ret < 0 {
						log.Trace("AvcodecReceivePacket err ", avutil.ErrorFromCode(ret))
						continue
					}
					if ret == 0 {
						data0 := pkt_enc.Data()
						buf := make([]byte, pkt_enc.GetPacketSize())
						start := uintptr(unsafe.Pointer(data0))
						for i := 0; i < pkt_enc.GetPacketSize(); i++ {
							elem := *(*uint8)(unsafe.Pointer(start + uintptr(i)))
							buf[i] = elem
						}
						
						file.Write(buf)
						//encLatency.End()
						fmt.Println(encLatency.End())
					}

					avutil.AvFrameUnref(frame)
					//avutil.AvFrameUnref(oframe)

				}

			}
		}
	}	
}

func main() {
	//decode_encode()
	decode_overlay_encode()
}