package main

import (
	"net"
	"io/ioutil"
	"time"
	"github.com/ailumiyana/streaming/rtp"
	"github.com/ailumiyana/streaming/h264"
	"os"
	"flag"
	"encoding/binary"
	"fmt"
)


func PackH264DataToNalus(bytes []byte) [][]byte {
	l := len(bytes)
	var startPos []int
	var nalus [][]byte
	j := 0 // split nalu in bytes to nalus 
	for i := 0; i < l - 5; i++ {
		if bytes[i] == 0 && bytes[i+1] == 0 && bytes[i+2] == 1 {
			if i > 0 && bytes[i-1] == 0 {//parameter set startpos
				startPos = append(startPos, i-1)
			} else {
				startPos = append(startPos, i)
			}
			j++
			if j > 1 {
				b := bytes[startPos[j-2]:startPos[j-1]]
				nalus = append(nalus, b)
			}
		}
	}
	nalus = append(nalus, bytes[startPos[j-1]:])
	if len(nalus) != len(startPos) {
		panic("unknown error at split nalu in bytes to nalus ")
	}

	return nalus
}

var audio_data []byte

func rtp_send_aac() {

	addAuHeader := func(b []byte) []byte {
		/* au header
	
		  +---------------------------------------+
	
		  |     AU-size                           |
	
		  +---------------------------------------+
	
		  |     AU-Index / AU-Index-delta         |
	
		  +---------------------------------------+
	
		  |     CTS-flag                          |
	
		  +---------------------------------------+
	
		  |     CTS-delta                         |
	
		  +---------------------------------------+
	
		  |     DTS-flag                          |
	
		  +---------------------------------------+
	
		  |     DTS-delta                         |
	
		  +---------------------------------------+
	
		  |     RAP-flag                          |
	
		  +---------------------------------------+
	
		  |     Stream-state                      |
	
		  +---------------------------------------+ */
	
		auHeader := []byte{0x00, 0x10}
		auHeaderLen := []byte{0x00, 0x00}
		auHeaderLen[0] = byte((uint16(len(b)) & 0x1fe0) >> 5)
		auHeaderLen[1] = byte((uint16(len(b)) & 0x001f) << 3)
	
		auHeader = append(auHeader, auHeaderLen[0:2]...)
	
		return append(auHeader, b[0:]...)
	}

	lo_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:" + *aport)
	if err != nil{
		fmt.Println("net ResolveUDPAddr Error :", err)
		os.Exit(1)
	}

	raddr, err := net.ResolveUDPAddr("udp", *aaddr)
	if err != nil{
		fmt.Println("net ResolveUDPAddr Error :", err)
	}
	
	fmt.Println("remote vedio addresses : ", raddr.IP, ":", raddr.Port)

	conn, err := net.DialUDP("udp4", lo_addr, raddr)
	if err != nil {
		fmt.Println("net DialUDP.")
		return 
	}

	defer conn.Close()

	b := make([]byte, 2)
	rp := rtp.NewDefaultPacketWith96Type()
	rp.SetSequence(50000)
	l := len(audio_data) 
	aRead := 0

	time_incr := time.Duration(0)
	
	for {
		aRead = 0
		for aRead < l {
			if audio_data[aRead] != 0xff && (audio_data[aRead+1] & 0xf0) != 0xf0 {
				fmt.Println("audio_data error")
			}

			b[0] = ((audio_data[aRead+3] & 0x03) << 3) | (( audio_data[aRead+4] & 0xe0) >> 5 )
			b[1] = (audio_data[aRead+4] << 3) | ((audio_data[aRead+5] & 0xe0) >> 5) 

			le := int(binary.BigEndian.Uint16(b))
			
			//fmt.Println(hex.EncodeToString(audio_data[aRead:aRead+le]))

			//add AU Header
			payload := addAuHeader(audio_data[aRead+7:aRead+le])

			// rtp send
			rp.SetPayload(payload)

			rp.SetTimeStamp(rp.TimeStamp() + 1024)
			//rp.SetTimeStamp(uint32(( time_wheel.Now() / time.Millisecond ) * 441 / 10))

			rp.SetSequence(rp.Sequence() + 1)
			
			//fmt.Println("seq : ", rp.Sequence(), "timestamp : ", rp.TimeStamp())

			//fmt.Println("----------------aac-:", time.Duration(rp.TimeStamp()) * 10 / 441 * time.Millisecond)
			time_incr += 23220*time.Microsecond

			//fmt.Println(hex.EncodeToString(rp.GetRtpBytes()))

			conn.Write(rp.GetRtpBytes())

			aRead += le

			time.Sleep(23220*time.Microsecond)

		}
	}
}

var aaddr *string
var aport *string

//e.g  :
//sudo go run h264_rtp_send.go -aip=192.168.0.78:1230 -rip=192.168.0.78 -port=1236
func main()  {
	vaddr := flag.String("vaddr", "127.0.0.1:6002", "remote video addr")
	aaddr = flag.String("aaddr", "127.0.0.1:6001", "remote audio addr ")
	lport := flag.String("lport", "10001", "local video port")
	aport  =  flag.String("aport", "10002", "local audio port")
	file_path := flag.String("file", "1080p.h264", "h264 file path")
	afile_path := flag.String("afile", "ilem.aac", "aac file path")
	flag.Parse()	

	raddr, err := net.ResolveUDPAddr("udp", *vaddr)
	if err != nil{
		fmt.Println("net ResolveUDPAddr Error : ", err)
		os.Exit(1)
	}

	laddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:" + *lport)
	if err != nil{
		fmt.Println("net ResolveUDPAddr Error : ", err)
		os.Exit(1)
	}
	
	fmt.Println("remote vedio addresses : ", raddr.IP, ":", raddr.Port)

	conn, err := net.DialUDP("udp4", laddr, raddr)
	if err != nil {
		fmt.Println("net DialUDP.")
		return 
	}

	defer conn.Close()

	data, err := ioutil.ReadFile(*file_path)
    if err != nil {
        fmt.Println("File reading error", err)
        return
	}
	fmt.Println("Open Success.")
	l := len(data)
    fmt.Println("size of file:", l)
	
	audio_data, err = ioutil.ReadFile(*afile_path)
    if err != nil {
        fmt.Println("File reading error", err)
        return
	}

	go rtp_send_aac()

	rtpPacket  := rtp.NewDefaultPacketWith97Type()
	rtpPacket.SetSequence(50000)

	nalus := PackH264DataToNalus(data)
	
	h264Praser := h264.NewParser()
	last_h264_tpye := h264.NalueTypeNotDefined
	time_incr := time.Duration(0)

	for {
		for _, v := range nalus {
			rps := rtpPacket.ParserNaluToRtpPayload(v)

			// H264 30FPS : 90000 / 30 : diff = 3000
			if last_h264_tpye == h264.NalueTypeSlice || last_h264_tpye == h264.NalueTypeIdr {
				rtpPacket.SetTimeStamp(rtpPacket.TimeStamp() + 3000)
				//rtpPacket.SetTimeStamp(uint32(( time_wheel.Now() / time.Millisecond ) * 90))
				time.Sleep(33334*time.Microsecond)
				time_incr += 33*time.Millisecond
			}

			h264Praser.FillNaluHead(rps[0][0])

			if h264Praser.NaluType() == h264.NalueTypeFuA {
				h264Praser.FillShadUnitA([2]byte{rps[0][0], rps[0][1]})
				last_h264_tpye = h264Praser.ShardA().NaluType()
			} else {
				last_h264_tpye = h264Praser.NaluType()
			}
			
			//fmt.Println("----------------h264-", last_h264_tpye, time.Duration(rtpPacket.TimeStamp() / 90) * time.Millisecond, rtpPacket.TimeStamp())
			//fmt.Println("----------------h264-", last_h264_tpye, time_incr, time.Duration(rtpPacket.TimeStamp() / 90) * time.Millisecond, rtpPacket.TimeStamp(), time.Duration(rtpPacket.TimeStamp()) * 1000 / 90000 * time.Millisecond )

			for _, q := range rps {
				rtpPacket.SetSequence(rtpPacket.Sequence() + 1)
				rtpPacket.SetPayload(q)
				conn.Write(rtpPacket.GetRtpBytes())
				//lconn.WriteToUDP(rtpPacket.GetRtpBytes(), &net.UDPAddr{IP: net.ParseIP("192.168.0.78"), Port: 1236})
			}
		}
	}
}
