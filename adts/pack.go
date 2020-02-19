package adts

import(
	//"fmt"
	"encoding/binary"
)

const (
	ADTS_HEADER_SIZE = 7
	ADTS_MAX_FRAME_BYTES = ((1 << 13) - 1)
)

type AdtsHeader [7]byte

func index(i int) int{
	return i/8
}

func Default() *AdtsHeader {
	bits := AdtsHeader{0,0,0,0,0,0,0}

	//sync(12 bits) 帧同步标识一个帧的开始，固定为0xFFF
	//sync := 0xfff0
	bits[index(0)] |= byte(0xff)
	bits[index(8)] |= byte(0xf0)

	//id(1 bits) 0:MPEG-4 / 1:MPEG-2
	//bits.Set(12)
	id := byte(0x00)
	bits[index(12)] |= id

	//layer(2 bits) 固定:'00'
	layer := byte(0x00)
	//bits.Clear(13)
	//bits.Clear(14)
	bits[index(13)] |= layer

	//protection_absent(1 bits) 标识是否进行误码校验。0表示有CRC校验，1表示没有CRC校验
	//bits.Set(15)
	protection_absent := byte(0x01)
	bits[index(15)] |= protection_absent

	//profile(2 bits)：标识使用哪个级别的AAC。0: AAC Main 1:AAC LC (Low Complexity) 
	//2:AAC SSR (Scalable Sample Rate) 3:AAC LTP (Long Term Prediction)
	//bits.Set(16)  // LC
	//bits.Clear(17)
	profile := byte(0x40)
	bits[index(16)] |= profile 

	//sampling_frequency_index(4 bits)：标识使用的采样率的下标
	// 4 : 44100 //00|01 00|00
	//bits.Clear(18)
	//bits.Set(19)
	//bits.Clear(20)
	//bits.Clear(21)
	sampling_frequency_index := byte(4)
	bits[index(18)] |= sampling_frequency_index << 2

	//private_bit(1 bits)：私有位，编码时设置为0，解码时忽略
	//bits.Clear(22)
	private_bit := byte(0x00)
	bits[index(22)] |= private_bit 

	//channel_configuration(3 bits)：标识声道数
	//default : 2 // 0 / 10  
	//bits.Clear(23)
	//bits.Set(24)
	//bits.Clear(25)
	chans := byte(2)
	bits[index(23)] |= chans >> 2
	bits[index(24)] |= chans << 6

	//original_copy(1 bits)：编码时设置为0，解码时忽略
	bits[index(26)] |= 0x00

	//home(1 bits)：编码时设置为0，解码时忽略
	bits[index(27)] |= 0x00


	// adts_variable_header
	copyright_identification_bit := byte(0x00)
	bits[index(28)] |= copyright_identification_bit

	copyright_identification_start := byte(0x00)
	bits[index(29)] |= copyright_identification_start

	//aac_frame_length
	full_frame_size := uint16(0)
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, full_frame_size)

	bits[index(30)] |= ((b[0] >> 3) & 0x03)
	bits[index(32)] |= ((b[0] << 5) | (b[1] >> 3))
	bits[index(40)] |= (b[1] << 5)

	//adts_buffer_fullness (11 bit) 0x7ff
	bits[index(43)] |= 0x1f
	bits[index(48)] |= 0xfc


	// number_of_raw_data_blocks_in_frame (2 bit) 00
	bits[index(48)] |= 0x00


	//fmt.Printf("%b\n", bits)

	//adts :
	//1111 1111 1111 0000 0110 0011 0010 0000‬
	return &bits
}

func (a *AdtsHeader) Bytes() []byte {
	return (*a)[0:]
}

func (a *AdtsHeader) SetFrameSize(len uint16) {
	full_frame_size := len + ADTS_HEADER_SIZE

	if 	full_frame_size > ADTS_MAX_FRAME_BYTES {
		//"ADTS frame size too large\n"
		return 
	}

	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, full_frame_size)

	a.Bytes()[index(30)] &= 0xfc
	a.Bytes()[index(32)] &= 0x00
	a.Bytes()[index(40)] &= 0x1f

	a.Bytes()[index(30)] |= ((b[0] >> 3) & 0x03)
	a.Bytes()[index(32)] |= ((b[0] << 5) | (b[1] >> 3))
	a.Bytes()[index(40)] |= (b[1] << 5)
}

func (a *AdtsHeader) SetSamplingFrequency(freq int) {
	var sampling_frequency_index byte
	switch freq {
	case 48000:
		sampling_frequency_index = 3
	default:
		sampling_frequency_index = 4 //44100
	}

	a.Bytes()[index(18)] &= 0xc3
	a.Bytes()[index(18)] |= sampling_frequency_index << 2

}

func (a *AdtsHeader) SetChannels(chans byte) {

	if chans > 7 || chans < 1 {
		return
	}
	
	a.Bytes()[index(23)] &= 0xfe
	a.Bytes()[index(24)] &= 0x3f

	a.Bytes()[index(23)] |= chans >> 2
	a.Bytes()[index(24)] |= chans << 6
}