package signaling

import (
	"fmt"
	"io"

	"github.com/pion/webrtc/v3"
)

func JoinRoom(channel_id string, client *Client) {
	existingChannels.mu.Lock()
	defer existingChannels.mu.Unlock()

	channel := existingChannels.channels[channel_id]

	if channel == nil {
		usersMap := make(map[string]*Client)
		channel = &Channel{
			channel_id:            channel_id,
			users:                 usersMap,
			remoteTrackForwarders: make(map[string]*trackForwarder),
		}
		existingChannels.channels[channel_id] = channel
	}
	err := channel.addUser(client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Room succesfuly joined")
}
func GetClients(channel_id string) {
	existingChannels.mu.Lock()
	defer existingChannels.mu.Unlock()

	channel := existingChannels.channels[channel_id]
	channel.mu.Lock()
	defer channel.mu.Unlock()
	for _, clientId := range channel.users {
		fmt.Println("connected: ", clientId)
	}
}
func NewTrackForwarder(track *webrtc.TrackRemote) *trackForwarder {
	//track forwarder object is needed to stream remoteTrack (track that client stream to server)
	//to other users. Track local staticRTP stands for one track that will be streamed to the user
	tf := &trackForwarder{
		remoteTrack: track,
		localTracks: []*webrtc.TrackLocalStaticRTP{},
	}

	go tf.run() //run goroutine to start stream immediately even if our local tracks is nill

	return tf
}

// run is the core forwarding loop for a single TrackForwarder.
// This function runs forever (until the track ends).
// It reads packets from the *one* remoteTrack and writes them
// to *all* of the localTracks (the viewers).
func (tf *trackForwarder) run() {
	rtpBuf := make([]byte, 1500)

	for {
		i, _, readErr := tf.remoteTrack.Read(rtpBuf)

		if readErr != nil {
			if readErr == io.EOF {
				// The track has ended. This is normal. We stop the loop.
				fmt.Printf("Track %s ended\n", tf.remoteTrack.ID())
				return
			}
			// A different error occurred
			fmt.Printf("Error reading from remote track %s: %v\n", tf.remoteTrack.ID(), readErr)
			return
		}

		tf.mu.RLock()
		for _, localTrack := range tf.localTracks {
			if _, writeErr := localTrack.Write(rtpBuf[:i]); writeErr != nil {

			}
		}
		tf.mu.RUnlock()
	}
}

func (tf *trackForwarder) AddSubscriber(PeerConnection *webrtc.PeerConnection) {
	tf.mu.Lock()
	defer tf.mu.Unlock()

	localTrack, err := webrtc.NewTrackLocalStaticRTP(
		tf.remoteTrack.Codec().RTPCodecCapability,
		tf.remoteTrack.ID(),
		tf.remoteTrack.StreamID(),
	)
	if err != nil {
		fmt.Println("error during creation of local track")
	}
	_, err = PeerConnection.AddTrack(localTrack)
	if err != nil {

	}
	/*
		go func() {
			rtcpBuf := make([]byte, 1500)
			for {
				n, _, rtcpErr := rtpSender.Read(rtcpBuf)
				if rtcpErr != nil {
					if rtcpErr == io.EOF {
						return
					}
					fmt.Println("Error reading RTCP1:", rtcpErr)
					return
				}
				if _, rtcpErr := tf.remoteTrack.write(rtcpBuf[:n]); rtcpErr != nil {
					fmt.Println("Error writing RTCP2:", rtcpErr)
				}
			}
		}()
	*/
	tf.localTracks = append(tf.localTracks, localTrack)
	fmt.Printf("Added subscriber. Track %s now has %d subscribers.\n",
		tf.remoteTrack.ID(), len(tf.localTracks))
}
