package main

import (
	log "github.com/astaxie/beego/logs"
	"github.com/ailumiyana/goav-incr/goav/avutil"
	"github.com/ailumiyana/goav-incr/goav/avcodec"
	"github.com/ailumiyana/goav-incr/goav/avfilter"

	"github.com/ailumiyana/tools/latency"

	"fmt"
	"time"
	"io/ioutil"
	"os"
	"strconv"
	"unsafe"
)

var input_h264  string = "1080p.h264"
var input_h264_720  string = "recv0.h264"
var output_h264 string = "out.h264"


var overlay16 string = "color=white:r=30:size=1920x1080:sar=1/1,hwupload=extra_hw_frames=300,format=qsv [background];" +
"[in0] scale_qsv=479:269 [in_0];" +
"[in1] scale_qsv=479:269 [in_1];" +
"[in2] scale_qsv=479:269 [in_2];" +
"[in3] scale_qsv=479:269 [in_3];" +
"[in4] scale_qsv=479:269 [in_4];" +
"[in5] scale_qsv=479:269 [in_5];" +
"[in6] scale_qsv=479:269 [in_6];" +
"[in7] scale_qsv=479:269 [in_7];" +
"[in8] scale_qsv=479:269 [in_8];" +
"[in9] scale_qsv=479:269 [in_9];" +
"[in10] scale_qsv=479:269 [in_10];" +
"[in11] scale_qsv=479:269 [in_11];" +
"[in12] scale_qsv=479:269 [in_12];" +
"[in13] scale_qsv=479:269 [in_13];" +
"[in14] scale_qsv=479:269 [in_14];" +
"[in15] scale_qsv=479:269 [in_15];" +
"[background][in_0] overlay_qsv=x=0:y=0 [background+in0_scale];" +
"[background+in0_scale][in_1] overlay_qsv=x=481:y=0 [background+in0_scale+in1_scale];" +
"[background+in0_scale+in1_scale][in_2] overlay_qsv=x=961:y=0 [background+in0_scale+in1_scale+in2_scale];" +
"[background+in0_scale+in1_scale+in2_scale][in_3] overlay_qsv=x=1441:y=0 [background+in0_scale+in1_scale+in2_scale+in3_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale][in_4] overlay_qsv=x=0:y=271 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale][in_5] overlay_qsv=x=481:y=271 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale][in_6] overlay_qsv=x=961:y=271 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale][in_7] overlay_qsv=x=1441:y=271 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale][in_8] overlay_qsv=x=0:y=541 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale][in_9] overlay_qsv=x=481:y=541 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale][in_10] overlay_qsv=x=961:y=541 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale+in10_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale+in10_scale][in_11] overlay_qsv=x=1441:y=541 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale+in10_scale+in11_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale+in10_scale+in11_scale][in_12] overlay_qsv=x=0:y=811 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale+in10_scale+in11_scale+in12_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale+in10_scale+in11_scale+in12_scale][in_13] overlay_qsv=x=481:y=811 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale+in10_scale+in11_scale+in12_scale+in13_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale+in10_scale+in11_scale+in12_scale+in13_scale][in_14] overlay_qsv=x=961:y=811 [background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale+in10_scale+in11_scale+in12_scale+in13_scale+in14_scale];" +
"[background+in0_scale+in1_scale+in2_scale+in3_scale+in4_scale+in5_scale+in6_scale+in7_scale+in8_scale+in9_scale+in10_scale+in11_scale+in12_scale+in13_scale+in14_scale][in_15] overlay_qsv=x=1441:y=811"


type QsvHWDeviceCtx struct {
	hw_device_ctx *avutil.BufferRef
}

// input : "/dev/dri/card0" or "/dev/dri/renderD128"
func Create(device string) *QsvHWDeviceCtx {
	var hw_device_ctx *avutil.BufferRef
	err := avutil.AVHwdeviceCtxCreate(&hw_device_ctx, avutil.AV_HWDEVICE_TYPE_QSV,
		device, nil, 0)
	if err < 0 {
		log.Error("AVHwdeviceCtxCreate err : ", avutil.ErrorFromCode(err))
		return nil
	}

	return &QsvHWDeviceCtx{
		hw_device_ctx,
	}
}

func (v *QsvHWDeviceCtx) Context() *avutil.BufferRef{
	return v.hw_device_ctx
}

//warning need unref by avutil.AVBufferUnref()
func (v *QsvHWDeviceCtx)GetAnAllocHwframeCtxRef(format, sw_format avutil.PixelFormat, w, h, s int) *avutil.BufferRef{
	hw_frames_ref := avutil.AVHwframeCtxAlloc(v.hw_device_ctx)
	if hw_frames_ref == nil {
		log.Error("AVHwframeCtxAlloc err")
		return nil
	}

	frames_ctx := hw_frames_ref.HWFramesContext()

	frames_ctx.SetQsvHWFramesContextPrarms(format, sw_format, w, h, s)

	err := avutil.AVHwFrameCtxInit(hw_frames_ref)
	if err < 0 {
		log.Error("AVHwFrameCtxInit err : ", avutil.ErrorFromCode(err))
		avutil.AVBufferUnref(&hw_frames_ref)
		return nil
	}

	return hw_frames_ref
}



func decode_overlay_encode() {
	//decoder
	qsv_device  := Create("/dev/dri/renderD128")

	pkt 		   := avcodec.AvPacketAlloc()
	if pkt == nil {
		log.Critical("AvPacketAlloc failed.")
		return 
	}

	codec 		 := avcodec.AvcodecFindDecoderByName("h264_qsv")
	if codec == nil {
		log.Critical("AvcodecFindDecoderByName failed.")
		return
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

	context.SetHWDeviceCtx(avutil.AVBufferRef(qsv_device.Context()))
	err := context.AvcodecOpen2(codec, nil)
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return 
	}
	context.SetDefaultQsvGetFormat()

	// encoder
	codec_enc 		   := avcodec.AvcodecFindEncoderByName("h264_qsv")
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
		avcodec.AV_PIX_FMT_QSV,
		false, 10)
	context_enc.SetTimebase(1, 30)

	hw_frames_ref := qsv_device.GetAnAllocHwframeCtxRef(
							avutil.PixelFormat(avcodec.AV_PIX_FMT_QSV),
							avutil.PixelFormat(avcodec.AV_PIX_FMT_NV12), 1920, 1080, 64)


	/*context_enc.SetHWFramesCtx(avutil.AVBufferRef(outs[0].AvBuffersinkGetHwFramesCtx()))

	err = context_enc.AvcodecOpen2(codec_enc, nil)
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return
	}*/

	decLatency := latency.New("vaapi", "decode")
	encLatency := latency.New("vaapi", "encode")

	// filter
	args := "video_size=1920x1080:pix_fmt=qsv:time_base=1/30:pixel_aspect=1/1"
	//des  := "[in0] scale_qsv=w=960:h=540 [s0];color=black:r=30:size=1920x1080:sar=1/1, hwupload=extra_hw_frames=64,format=qsv [b0];[b0][s0] overlay_qsv=x=0:y=0"
	des  := overlay16

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

	if qsv_device.Context() != nil {
		for _, v := range graph.Filters() {
			v.SetHWDeviceCtx(avutil.AVBufferRef(qsv_device.Context()))
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

	overlayLatency := latency.New("qsv", "overlay")

	context_enc.SetHWFramesCtx(avutil.AVBufferRef(outs[0].AvBuffersinkGetHwFramesCtx()))

	var dict *avutil.Dictionary = avutil.AvDictAlloc()
	er := dict.AvDictSet("profile", "high", 0)
	if er < 0 {
		log.Critical("AvDictSet failed")
		return 
	}
	er = dict.AvDictSet("level", "52", 0)
	if er < 0 {
		log.Critical("AvDictSet failed")
		return 
	}

	fmt.Println(dict.AvDictCount())


	err = context_enc.AvcodecOpen2(codec_enc, (**avcodec.Dictionary)(unsafe.Pointer(&dict)))
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return
	}
	fmt.Println(dict.AvDictCount())

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
				
				var fs []*avutil.Frame

				for i := 0 ; i < 16 ; i++ {
					fs = append(fs, avutil.AvFrameClone(frame))
					ret = avfilter.AvBuffersrcAddFrame(ins[i], (*avfilter.Frame)(unsafe.Pointer(fs[i])))
					if ret < 0 {
						for i := 0 ; i < len(fs) ; i++ {
							avutil.AvFrameFree(fs[i])
						}
						fmt.Println("AvBuffersrcAddFrame error,", avutil.ErrorFromCode(ret))
						continue
					}
				}
				avutil.AvFrameUnref(frame)
				/*ret = avfilter.AvBuffersrcAddFrame(ins[0], (*avfilter.Frame)(unsafe.Pointer(frame)))
				if ret < 0 {
					fmt.Println("AvBuffersrcAddFrame error,", avutil.ErrorFromCode(ret))
					continue
				}*/

				ret = avfilter.AvBufferSinkGetFrame(outs[0], (*avfilter.Frame)(unsafe.Pointer(frames[0])))
				for i := 0 ; i < len(fs) ; i++ {
					avutil.AvFrameFree(fs[i])
				}
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



func filter_inpute_args() {
	//decoder
	qsv_device  := Create("/dev/dri/renderD128")

	pkt 		   := avcodec.AvPacketAlloc()
	if pkt == nil {
		log.Critical("AvPacketAlloc failed.")
		return 
	}

	codec 		 := avcodec.AvcodecFindDecoderByName("h264_qsv")
	if codec == nil {
		log.Critical("AvcodecFindDecoderByName failed.")
		return
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

	context.SetHWDeviceCtx(avutil.AVBufferRef(qsv_device.Context()))
	err := context.AvcodecOpen2(codec, nil)
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return 
	}
	context.SetDefaultQsvGetFormat()

	// encoder
	codec_enc 		   := avcodec.AvcodecFindEncoderByName("h264_qsv")
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
		avcodec.AV_PIX_FMT_QSV,
		false, 10)
	context_enc.SetTimebase(1, 30)

	hw_frames_ref := qsv_device.GetAnAllocHwframeCtxRef(
								avutil.PixelFormat(avcodec.AV_PIX_FMT_QSV),
								avutil.PixelFormat(avcodec.AV_PIX_FMT_NV12), 1920, 1080, 64)

	pkt1 		   := avcodec.AvPacketAlloc()
	if pkt == nil {
		log.Critical("AvPacketAlloc failed.")
		return 
	}

	codec1 		 := avcodec.AvcodecFindDecoderByName("h264_qsv")
	if codec == nil {
		log.Critical("AvcodecFindDecoderByName failed.")
		return
	}

	context1 	   := codec.AvcodecAllocContext3()
	if context1 == nil {
		log.Critical("AvcodecAllocContext3 failed.")
		return 
	}

	parserContext1  := avcodec.AvParserInit(int(avcodec.CodecId(avcodec.AV_CODEC_ID_H264)))
	if parserContext1 == nil {
		log.Critical("AvParserInit failed.")
		return 
	}

	frame1   	   := avutil.AvFrameAlloc()
	if frame1 == nil {
		log.Critical("AvFrameAlloc failed.")
		return 
	}

	context1.SetHWDeviceCtx(avutil.AVBufferRef(qsv_device.Context()))
	err = context1.AvcodecOpen2(codec1, nil)
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return 
	}
	context1.SetDefaultQsvGetFormat()


	hw_frames_ref1 := qsv_device.GetAnAllocHwframeCtxRef(
							avutil.PixelFormat(avcodec.AV_PIX_FMT_QSV),
							avutil.PixelFormat(avcodec.AV_PIX_FMT_NV12), 1280, 720, 64)
	/*context_enc.SetHWFramesCtx(avutil.AVBufferRef(outs[0].AvBuffersinkGetHwFramesCtx()))

	err = context_enc.AvcodecOpen2(codec_enc, nil)
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return
	}*/

	decLatency := latency.New("vaapi", "decode")
	encLatency := latency.New("vaapi", "encode")

	// filter
	args := "video_size=1920x1080:pix_fmt=qsv:time_base=1/30:pixel_aspect=1/1"
	//des  := "[in0] scale_qsv=w=960:h=540 [s0];color=black:r=30:size=1920x1080:sar=1/1, hwupload=extra_hw_frames=64,format=qsv [b0];[b0][s0] overlay_qsv=x=0:y=0"
	des  := overlay16

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

	if qsv_device.Context() != nil {
		for _, v := range graph.Filters() {
			v.SetHWDeviceCtx(avutil.AVBufferRef(qsv_device.Context()))
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

	overlayLatency := latency.New("qsv", "overlay")

	context_enc.SetHWFramesCtx(avutil.AVBufferRef(outs[0].AvBuffersinkGetHwFramesCtx()))

	var dict *avutil.Dictionary = avutil.AvDictAlloc()
	er := dict.AvDictSet("profile", "high", 0)
	if er < 0 {
		log.Critical("AvDictSet failed")
		return 
	}
	er = dict.AvDictSet("level", "52", 0)
	if er < 0 {
		log.Critical("AvDictSet failed")
		return 
	}

	fmt.Println(dict.AvDictCount())


	err = context_enc.AvcodecOpen2(codec_enc, (**avcodec.Dictionary)(unsafe.Pointer(&dict)))
	if err < 0 {
		log.Critical("AvcodecOpen2 failed.")
		return
	}
	fmt.Println(dict.AvDictCount())

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
				
				var fs []*avutil.Frame

				for i := 0 ; i < 16 ; i++ {
					fs = append(fs, avutil.AvFrameClone(frame))
					ret = avfilter.AvBuffersrcAddFrame(ins[i], (*avfilter.Frame)(unsafe.Pointer(fs[i])))
					if ret < 0 {
						for i := 0 ; i < len(fs) ; i++ {
							avutil.AvFrameFree(fs[i])
						}
						fmt.Println("AvBuffersrcAddFrame error,", avutil.ErrorFromCode(ret))
						continue
					}
				}
				avutil.AvFrameUnref(frame)
				/*ret = avfilter.AvBuffersrcAddFrame(ins[0], (*avfilter.Frame)(unsafe.Pointer(frame)))
				if ret < 0 {
					fmt.Println("AvBuffersrcAddFrame error,", avutil.ErrorFromCode(ret))
					continue
				}*/

				ret = avfilter.AvBufferSinkGetFrame(outs[0], (*avfilter.Frame)(unsafe.Pointer(frames[0])))
				for i := 0 ; i < len(fs) ; i++ {
					avutil.AvFrameFree(fs[i])
				}
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
	avutil.AvLogSetLevel(avutil.AV_LOG_TRACE)
	//fmt.Println(filter.P1080InTestOverlay16)
	//fmt.Println(avcodec.AV_PIX_FMT_QSV)
	//fmt.Println(avcodec.AV_PIX_FMT_NV12)
	//fmt.Println(avutil.ErrorFromCode(-15))

	decode_overlay_encode()

	time.Sleep(time.Hour)
	//oclOpen()
}