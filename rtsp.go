package gortsplib

import (
	"fmt"
	"github.com/aliveyun/gortsplib/pkg/url"
	"log"
	//"os"
	"time"
	"bytes"
	"github.com/aliveyun/gortsplib/pkg/av"
	"github.com/aliveyun/gortsplib/pkg/formats"
	"github.com/aliveyun/gortsplib/pkg/media"
	"github.com/bluenviron/mediacommon/pkg/codecs/h264"
	"github.com/bluenviron/mediacommon/pkg/codecs/h265"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
)

// Packet stores compressed audio/video data.
type Packet struct {
	IsKeyFrame      bool // video packet is key frame
	Idx             int8 // stream index in container format
	VideoCodec      av.CodecType
	CompositionTime time.Duration // packet presentation time minus decode time for H264 B-Frame
	Time            uint32
	Duration        time.Duration //packet duration
	Data            []byte        // packet data
}
type RTSPClientOptions struct {
	Debug            bool
	URL              string
	DialTimeout      time.Duration
	ReadWriteTimeout time.Duration
	DisableAudio     bool
	OutgoingProxy    bool
}

type AVDecode interface {
	Decode(pkt *rtp.Packet) ([][]byte, time.Duration, error)
}
type AUDecode interface {
	Decode(pkt *rtp.Packet) ([]byte, time.Duration, error)
}

type RTSPClient struct {
	medias       media.Medias
	audioTrackId int
	videoTrackId int
	//waitSend bool
	OutgoingPacketQueue chan *Packet
	Signals             chan int
	//audioCodec          av.CodecType
	//videoCodec          av.CodecType
	//sps []byte
	//pps []byte
	//track              interface{}
	//file *os.File
	Av AVDecode
	Au AUDecode
	//Tracks  gortsplib.Tracks
	//Avtrack  interface{}
	//Autrack  interface{}
	AvFormat interface{}
	AuFormat interface{}
	conn     *Client
	//OnPacketRTP func(*ClientOnPacketRTPCtx)
	// called when receiving a RTCP packet.
	//OnPacketRTCP func(*ClientOnPacketRTCPCtx)
}

func Dial(options RTSPClientOptions) (*RTSPClient, error) {

	// Client allows to set additional client options
	client := &RTSPClient{
		audioTrackId: -1,
		videoTrackId: -1,
		//waitSend: true,
		Signals:             make(chan int, 100),
		OutgoingPacketQueue: make(chan *Packet, 3000),
	}

	client.conn = &Client{
		// transport protocol (UDP, Multicast or TCP). If nil, it is chosen automatically
		Transport: nil,
		// timeout of read operations
		ReadTimeout: 10 * time.Second,
		// timeout of write operations
		WriteTimeout: 10 * time.Second,
	}
	// parse URL
	//u, err := url.Parse("rtsp://admin:admin@192.168.138.120:554/h264/ch1/main/av_stream") //265
	//u, err := url.Parse("rtsp://admin:abc12345@10.55.17.23:554/Streaming/Channels/1") //1
	//u, err := url.Parse("rtsp://admin:infore@123@10.55.128.71:554/Streaming/Channels/1")
	u, err := url.Parse(options.URL)
	if err != nil {
		return nil, err
	}
	// connect to the server
	err = client.conn.Start(u.Scheme, u.Host)
	if err != nil {
		return nil, err
	}
	//defer 	client.C.Close()
	//var baseURL  *url.URL
	// find published medias
	medias, baseURL, _, err := client.conn.Describe(u)
	if err != nil {
		return nil, err
	}
	//fmt.Println("kkkkk",medias)
	client.medias = medias

	// for _, track := range tracks {

	// 	log.Println("multiple video tracks are not supported", track)
	// 	switch ttrack := track.(type) {
	// 	case *TrackH265:
	// 		//track:=track.(*gortsplib.TrackH265)
	// 		//rtpDec : track.CreateDecoder()
	// 		log.Println("multiple video tracks are not supported", ttrack)
	// 	case *TrackH264:
	// 		//track:=track.(*gortsplib.TrackH264)
	// 		//rtpDec : track.CreateDecoder()
	// 		log.Println("multiple video tracks are not supported", ttrack)

	// 	case *TrackMPEG4Audio:
	// 		log.Println("multiple video tracks are not supported", ttrack)
	// 	}
	// }

	client.SetTracks()
	// setup and read all tracks
	//err = client.conn.SetupAndPlay(tracks, baseURL)
	err = client.conn.SetupAll(medias, baseURL)
	if err != nil {
		return nil, err
	}

	client.conn.OnPacketRTPAny(func(medi *media.Media, forma formats.Format, pkt *rtp.Packet) {

		client.RTPDemuxer(medi, forma, pkt)
	})

	// called when a RTCP packet arrives
	client.conn.OnPacketRTCPAny(func(medi *media.Media, pkt rtcp.Packet) {
		log.Printf("RTCP packet from media %v, type %T\n", medi, pkt)
	})
	_, err = client.conn.Play(nil)
	if err != nil {
		return nil, err
	}
	// go 	client.C.Wait()
	return client, nil
}

func (client *RTSPClient) Close() {
	if client.conn != nil {
		client.conn.Close()
	}
}

//var f *os.File

func (p *RTSPClient) RTPDemuxer(medi *media.Media, forma formats.Format, pkt *rtp.Packet) ([]*Packet, bool) {

	// if f == nil {
	// 	f, _ = os.Create("444444.h265")
	// 	fmt.Println("ldkldflklfglkfglkg")
	// }

	if "audio" == medi.Type {
		op, _, err := p.Au.Decode(pkt)
		if err != nil {
			log.Println("multiple video tracks are not supported", medi.Type)
			return nil, false
		}
		var retmap Packet
		retmap.Idx = int8(p.audioTrackId)
		switch forma.(type) {
		case *formats.G711:
			retmap.VideoCodec = av.PCM_ALAW
			retmap.Data = append(retmap.Data, op...)
		}
		retmap.Time = pkt.Timestamp * 1000 / 8000
		//log.Println("multiple audio tracks  retmap.Time",retmap.Time)
		p.OutgoingPacketQueue <- &retmap

	} else if "video" == medi.Type {

		flag := false
		var retmap Packet
		nalus, _, err := p.Av.Decode(pkt)
		if err != nil {
			//log.Printf("multiple video tracks are not supported flag=%d TrackID=%d ", flag, medi.Type)
			return nil, false
		}
		switch forma.(type) {
		case *formats.H265:
			retmap.VideoCodec = av.H265
			nalus, flag = p.h265RemuxNALUs(nalus, p.AvFormat.(*formats.H265))
		case *formats.H264:
			retmap.VideoCodec = av.H264
			nalus, flag = p.h264RemuxNALUs(nalus, p.AvFormat.(*formats.H264))
		}
		//nalus,_ = p.h264RemuxNALUs(nalus,p.Avtrack.(*gortsplib.TrackH264))

		retmap.IsKeyFrame = flag
		for _, nalu := range nalus {

			/*typ := h265.NALUType((nalu[0] & 0x7E) >> 1)
			switch typ {
			case h265.NALUType_PREFIX_SEI_NUT, h265.NALUType_SUFFIX_SEI_NUT:
				// remove since they're automatically added
				continue
			}*/
			if _,ok:=forma.(*formats.H264);ok{
				typ := h264.NALUType(nalu[0] & 0x1F)
				switch typ {
				case h264.NALUTypeSEI:
					// remove since they're automatically added
					continue
				}

			}else if _,ok:=forma.(*formats.H265);ok {
				typ := h265.NALUType((nalu[0] & 0x7E) >> 1)
				switch typ {
				case h265.NALUType_PREFIX_SEI_NUT, h265.NALUType_SUFFIX_SEI_NUT:
					// remove since they're automatically added
					continue
				}	
			}
		

			nalu = append([]byte{0, 0, 0, 1}, bytes.Join([][]byte{nalu[0:]}, []byte{0, 0, 0, 1})...)
			retmap.Data = append(retmap.Data, nalu...)
		}
		retmap.Time = pkt.Timestamp / 90

		//log.Println("multiple video tracks  retmap.Time",retmap.Time)
		if len(retmap.Data) > 0 {
			retmap.Idx = int8(p.videoTrackId)
			p.OutgoingPacketQueue <- &retmap
			//_, _ = f.Write(retmap.Data)
			//fmt.Println("lflfllf",len(retmap.Data))
			//f.Close()
			//panic("kkk")
		}

	}
	return nil, false
}

// remux is needed to fix corrupted streams and make streams
// compatible with all protocols.
func (t *RTSPClient) h264RemuxNALUsEx(nalus [][]byte, format *formats.H264) (filteredNALUs [][]byte, flag bool) {
	addSPSPPS := false
	n := len(nalus)
	for _, nalu := range nalus {
		typ := h264.NALUType(nalu[0] & 0x1F)
		if typ == h264.NALUTypeSPS || typ == h264.NALUTypePPS || typ == h264.NALUTypeIDR {
			if typ >= 1 && 5 >= typ {
				addSPSPPS = true
				n += 2
			}
			flag = true
			break

		}

	}

	if n == 0 {
		return nil, flag
	}

	filteredNALUs = make([][]byte, n)
	i := 0

	if addSPSPPS {
		filteredNALUs[0], filteredNALUs[1] = format.SafeParams()
		//filteredNALUs[0] = track.SafeSPS()
		//filteredNALUs[1] = track.SafePPS()
		i = 2
	}
	for _, nalu := range nalus {
		filteredNALUs[i] = nalu
		i++
	}
	// for _, nalu := range nalus {
	// 	typ := h264.NALUType(nalu[0] & 0x1F)
	// 	switch typ {
	// 	case h264.NALUTypeSPS, h264.NALUTypePPS:
	// 		// remove since they're automatically added
	// 		continue

	// 	case h264.NALUTypeAccessUnitDelimiter:
	// 		// remove since it is not needed
	// 		continue
	// 	}

	// 	filteredNALUs[i] = nalu
	// 	i++
	// }

	return filteredNALUs, flag
}
func (t *RTSPClient) h264RemuxNALUs(nalus [][]byte, format *formats.H264) (filteredNALUs [][]byte, flag bool) {
	addSPSPPS := false
	n := 0
	//havePPS := false
	for _, nalu := range nalus {
		typ := h264.NALUType(nalu[0] & 0x1F) //typ >= 1 && 5 >= typ
		if typ == h264.NALUTypeSPS || typ == h264.NALUTypePPS {
			//havePPS = true
			continue
		} else if typ == h264.NALUTypeIDR  {
			addSPSPPS = true
			n += 2
			flag = addSPSPPS	
		}
		n++
		
	}

	if n == 0 {
		return nil, flag
	}
	filteredNALUs = make([][]byte, n)
	i := 0
   
	if addSPSPPS {
		filteredNALUs[0], filteredNALUs[1] = format.SafeParams()
		i = 2
	}

	for _, nalu := range nalus {
		typ := h264.NALUType(nalu[0] & 0x1F)
		switch typ {
		case h264.NALUTypeSPS,h264.NALUTypePPS:
			// remove since they're automatically added
			continue

		}

		filteredNALUs[i] = nalu
		i++
	}

	return filteredNALUs, flag
}

// (code & 0x7E)>>1
func (t *RTSPClient) h265RemuxNALUs(nalus [][]byte, format *formats.H265) (filteredNALUs [][]byte, flag bool) {
	addSPSPPS := false
	n := 0
	//havePPS := false
	for _, nalu := range nalus {
		typ := h265.NALUType((nalu[0] & 0x7E) >> 1)
		//typ := h265.NALUType((nalu[0] & 0x7E) >> 1) // typ == h265.NALUType_SPS_NUT || typ == h265.NALUType_PPS_NUT ||  (typ >= 16 && 23 >= typ)
		if typ == h265.NALUType_VPS_NUT || typ == h265.NALUType_SPS_NUT || typ == h265.NALUType_PPS_NUT {
			//havePPS = true
			continue
		} else if typ >= 16 && typ >= 23 {
			addSPSPPS = true
			n += 3
			flag = addSPSPPS
		}
		n++
	}

	if n == 0 {
		return nil, flag
	}
	filteredNALUs = make([][]byte, n)
	i := 0

	if addSPSPPS {
		filteredNALUs[0], filteredNALUs[1], filteredNALUs[2] = format.SafeParams()
		i = 3
	}

	for _, nalu := range nalus {
		typ := h265.NALUType((nalu[0] & 0x7E) >> 1)
		switch typ {
		case h265.NALUType_VPS_NUT, h265.NALUType_SPS_NUT, h265.NALUType_PPS_NUT:
			// remove since they're automatically added
			continue

		}

		filteredNALUs[i] = nalu
		i++
	}

	return filteredNALUs, flag
}

// Decode [] AVDecode `json:"-"`
func (p *RTSPClient) SetTracks() error {
	for _, track := range p.medias {
		for trackId, forma := range track.Formats {
			fmt.Println("formats.", forma)
			switch forma.(type) {
			case *formats.H264:
				if p.Av == nil {
					var forma *formats.H264
					medi := p.medias.FindFormat(&forma)
					if medi == nil {
						log.Printf("media not found %s \n", "h264")
					}
					p.AvFormat = forma
					rtpDec := forma.CreateDecoder()
					p.Av = rtpDec
					p.videoTrackId = trackId
					fmt.Println("formats.H264", rtpDec)
				}
			case *formats.H265:
				if p.Av == nil {
					var forma *formats.H265
					medi := p.medias.FindFormat(&forma)
					if medi == nil {
						log.Printf("media not found %s \n", "h265")
					}
					p.AvFormat = forma
					rtpDec := forma.CreateDecoder()
					p.Av = rtpDec
					p.videoTrackId = trackId
				}
			case *formats.G711:
				if p.Au == nil {
					var forma *formats.G711

					medi := p.medias.FindFormat(&forma)
					if medi == nil {
						log.Printf("media not found %s \n", "g711")
					}
					p.AuFormat = forma
					rtpDec := forma.CreateDecoder()
					p.Au = rtpDec
					p.audioTrackId = trackId
				}
			case *formats.MPEG4Audio:
				if p.Au == nil {

				}
			}

		}
	}

	return nil

}
