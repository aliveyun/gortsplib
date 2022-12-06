package gortsplib

import (
	"errors"
	"log"
	"time"

	"github.com/aliveyun/gortsplib/pkg/url"
	//"github.com/golang/protobuf/ptypes/any"
	"bytes"
	"fmt"
	"github.com/aliveyun/gortsplib/pkg/av"
	"github.com/aliveyun/gortsplib/pkg/h264"
	"github.com/aliveyun/gortsplib/pkg/h265"
	//"os"
	"github.com/pion/rtp"
	"strings"
)

// Packet stores compressed audio/video data.
type Packet struct {
	IsKeyFrame      bool // video packet is key frame
	Idx             int8 // stream index in container format
	VideoCodec      av.CodecType
	CompositionTime time.Duration // packet presentation time minus decode time for H264 B-Frame
	Time            time.Duration // packet decode time
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
	tracks              Tracks
	audioTrackId        int
	videoTrackId        int
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
	Avtrack interface{}
	Autrack interface{}
	conn    *Client
	//OnPacketRTP func(*ClientOnPacketRTPCtx)
	// called when receiving a RTCP packet.
	//OnPacketRTCP func(*ClientOnPacketRTCPCtx)
}

func Dial(options RTSPClientOptions) (*RTSPClient, error) {
	client := &RTSPClient{
		audioTrackId:        -1,
		videoTrackId:        -1,
		Signals:             make(chan int, 100),
		OutgoingPacketQueue: make(chan *Packet, 3000),
	}
	//client.file, _ = os.Create("D://7777168.h264")
	client.conn = &Client{
		// called when a RTP packet arrives
		OnPacketRTP: func(ctx *ClientOnPacketRTPCtx) {
			//log.Printf("RTP packet from track %d, payload type %d\n", ctx.TrackID, ctx.Packet.Header.PayloadType)
			client.RTPDemuxer(ctx.TrackID, ctx.Packet)
		},
		// called when a RTCP packet arrives
		OnPacketRTCP: func(ctx *ClientOnPacketRTCPCtx) {
			log.Printf("RTCP packet from track %d, type %T\n", ctx.TrackID, ctx.Packet)
		},
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
	// find published tracks
	tracks, baseURL, _, err := client.conn.Describe(u)
	if err != nil {
		return nil, err
	}
	client.tracks = tracks

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
	err = client.conn.SetupAndPlay(tracks, baseURL)
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

func (p *RTSPClient) RTPDemuxer(TrackID int, pkt *rtp.Packet) ([]*Packet, bool) {

	if TrackID == p.audioTrackId {

		op, _, err := p.Au.Decode(pkt)
		if err != nil {
			log.Println("multiple video tracks are not supported")
			return nil, false
		}

		var retmap Packet
		switch p.Autrack.(type) {
		case *TrackG711:
			retmap.VideoCodec = av.PCM_ALAW
			retmap.Data = append(retmap.Data, op...)
		}
		p.OutgoingPacketQueue <- &retmap

	} else {
		return nil, false
		var flag bool
		var retmap Packet
		nalus, _, err := p.Av.Decode(pkt)
		if err != nil {
			return nil, false
			log.Println("multiple video tracks are not supported", nalus, flag)
		}
		switch p.Avtrack.(type) {
		case *TrackH265:
			retmap.VideoCodec = av.H265
			nalus, flag = p.h265RemuxNALUs(nalus, p.Avtrack.(*TrackH265))
		case *TrackH264:
			retmap.VideoCodec = av.H264
			nalus, flag = p.h264RemuxNALUs(nalus, p.Avtrack.(*TrackH264))
		}
		//nalus,_ = p.h264RemuxNALUs(nalus,p.Avtrack.(*gortsplib.TrackH264))

		retmap.IsKeyFrame = flag
		for _, nalu := range nalus {
			nalu = append([]byte{0, 0, 0, 1}, bytes.Join([][]byte{nalu[0:]}, []byte{0, 0, 0, 1})...)
			retmap.Data = append(retmap.Data, nalu...)
		}
		//p.file.Write(retmap.Data)
		p.OutgoingPacketQueue <- &retmap

		//log.Println("multiple video tracks are not supported")
	}
	return nil, false
}

// remux is needed to fix corrupted streams and make streams
// compatible with all protocols.
func (t *RTSPClient) h264RemuxNALUs(nalus [][]byte, track *TrackH264) (filteredNALUs [][]byte, flag bool) {
	addSPSPPS := false
	n := 0
	for _, nalu := range nalus {
		typ := h264.NALUType(nalu[0] & 0x1F)
		switch typ {
		case h264.NALUTypeSPS, h264.NALUTypePPS:
			continue
		case h264.NALUTypeAccessUnitDelimiter:
			continue
		case h264.NALUTypeIDR:
			// prepend SPS and PPS to the group if there's at least an IDR
			if !addSPSPPS {
				addSPSPPS = true
				n += 2
				flag = addSPSPPS
			}
		}
		n++
	}

	if n == 0 {
		return nil, flag
	}

	filteredNALUs = make([][]byte, n)
	i := 0

	if addSPSPPS {
		filteredNALUs[0] = track.SafeSPS()
		filteredNALUs[1] = track.SafePPS()
		i = 2
	}

	for _, nalu := range nalus {
		typ := h264.NALUType(nalu[0] & 0x1F)
		switch typ {
		case h264.NALUTypeSPS, h264.NALUTypePPS:
			// remove since they're automatically added
			continue

		case h264.NALUTypeAccessUnitDelimiter:
			// remove since it is not needed
			continue
		}

		filteredNALUs[i] = nalu
		i++
	}

	return filteredNALUs, flag
}

// (code & 0x7E)>>1
func (t *RTSPClient) h265RemuxNALUs(nalus [][]byte, track *TrackH265) (filteredNALUs [][]byte, flag bool) {
	addSPSPPS := false
	n := 0
	for _, nalu := range nalus {
		typ := h265.NALUType((nalu[0] & 0x7E) >> 1)
		if typ == h265.NAL_SPS || typ == h265.NAL_PPS {
			continue
		} else if typ >= 16 && typ >= 23 {
			if !addSPSPPS {
				addSPSPPS = true
				n += 2
				flag = addSPSPPS
			}
		}
		n++
	}
	if n == 0 {
		return nil, flag
	}
	filteredNALUs = make([][]byte, n)
	i := 0

	if addSPSPPS {
		filteredNALUs[0] = track.SafeSPS()
		filteredNALUs[1] = track.SafePPS()
		i = 2
	}

	for _, nalu := range nalus {
		typ := h265.NALUType((nalu[0] & 0x7E) >> 1)
		switch typ {
		case h265.NAL_SPS, h265.NAL_PPS:
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
	//p.Decode = make([]AVDecode, len(p.tracks))
	for trackId, track := range p.tracks {
		md := track.MediaDescription()
		v, ok := md.Attribute("rtpmap")
		if !ok {
			return errors.New("rtpmap attribute not found")
		}
		v = strings.TrimSpace(v)
		vals := strings.Split(v, " ")
		if len(vals) != 2 {
			continue
		}
		fmtp := make(map[string]string)
		if v, ok = md.Attribute("fmtp"); ok {
			if tmp := strings.SplitN(v, " ", 2); len(tmp) == 2 {
				for _, kv := range strings.Split(tmp[1], ";") {
					kv = strings.Trim(kv, " ")

					if len(kv) == 0 {
						continue
					}
					tmp := strings.SplitN(kv, "=", 2)
					if len(tmp) == 2 {
						fmtp[strings.TrimSpace(tmp[0])] = tmp[1]
					}
				}
			}
		}
		//timeScale := 0
		keyval := strings.Split(vals[1], "/")
		// if i, err := strconv.Atoi(keyval[1]); err == nil {
		// 	//timeScale = i
		// }
		if len(keyval) >= 2 {

			log.Printf("RTCP pack type %s \n", keyval[0])
			switch strings.ToLower(keyval[0]) {
			case "h264":
				p.videoTrackId = trackId
				track := track.(*TrackH264)
				p.Avtrack = track
				p.Av = track.CreateDecoder()
				log.Printf("RTCP packet from track  type %s %d\n", "h264", trackId)
			case "h265", "hevc":
				p.videoTrackId = trackId
				track := track.(*TrackH265)
				p.Avtrack = track
				p.Av = track.CreateDecoder()
				log.Printf("RTCP packet from track  type %s %d\n", "h265", trackId)
			case "pcma":
				p.audioTrackId = trackId
				track := track.(*TrackG711)
				p.Autrack = track
				p.Au = track.CreateDecoder()
				log.Printf("RTCP packet from track  type %s %d\n", "pcma", trackId)
			case "pcmu":
				p.audioTrackId = trackId
				track := track.(*TrackG711)
				p.Autrack = track
				p.Au = track.CreateDecoder()
				log.Printf("RTCP packet from track  type %s  %d\n", "pcmu", trackId)
			case "mpeg4-generic":
				p.audioTrackId = trackId
				log.Printf("RTCP packet from track  type %s %d\n", "generic", trackId)
			default:
				return fmt.Errorf("unsupport codec:%s", keyval[0])
			}
		}
	}
	return nil
}
