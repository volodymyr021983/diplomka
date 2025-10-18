<template>
    <div id="video-grid" ref="videoGrid">
        <!-- Our local video will always be here -->
        <video ref="localVideoPlayer" autoplay playsinline muted></video>
    </div>
    <button @click="ConnVoice">Connect to Voice</button>
</template>

<script setup lang="js">
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const localVideoPlayer = ref(null);
const videoGrid = ref(null); // A ref to the main video container
const route = useRoute();

// --- State Management ---
let publisherPC;
const subscriberPCs = {};
let myPeerID;

// Assuming serverId and channelId are available from the route
const serverId = route.params.server_id || 'test-server';
const channelId = route.params.channel_id || 'test-channel';

// Standard WebRTC configuration
let peerConfiguration = {
    iceServers: [{
        urls: [
            'stun:stun.l.google.com:19302',
            'stun:stun1.l.google.com:19302'
        ]
    }]
};

// --- WebSocket Connection ---
// This connects to the Go server you have in the Canvas.
const ws = new WebSocket(`wss://192.168.44.254:443/api/server/voice/${serverId}/${channelId}`);

// Assign the single message handler
ws.onmessage = OnMessageHandler;

// --- Main Message Handler ---
// This function acts as the central router for all messages from the server.
async function OnMessageHandler(event) {
    const msg = JSON.parse(event.data);
    console.log("<- Received message:", msg.type, msg.payload);

    switch (msg.type) {
        // STEP 2: The server confirms we have joined and gives us our unique ID.
        case 'room_joined':
            myPeerID = msg.payload.my_peer_id;
            console.log(`Successfully joined room. My Peer ID is ${myPeerID}`);
            // STEP 3: Now that we have our ID, we can start publishing our own media.
            startPublishing();
            break;
        
        case 'answer':
            // This handles answers for both our publisher and any subscriber connections.
            const answerPC = (msg.payload.remote_peer_id === myPeerID)
                ? publisherPC
                : subscriberPCs[msg.payload.remote_peer_id];

            if (answerPC) {
                await answerPC.setRemoteDescription(new RTCSessionDescription(msg.payload.sdp));
            }
            break;

        case 'offer':
            // The server is offering us a stream from an existing peer.
            SubscriptionOffer(msg.payload);
            break;

        case 'candidate':
            // The server sent an ICE candidate for one of our connections.
            RemoteCandidate(msg.payload);
            break;
    }
}

// --- User Actions ---

// STEP 1: The user clicks the button to start the connection process.
function ConnVoice() {
    // We need to wait for the WebSocket to be open before sending messages.
    if (ws.readyState === WebSocket.OPEN) {
        sendJoinRequest();
    } else {
        ws.onopen = sendJoinRequest;
    }
}

function sendJoinRequest() {
    console.log("-> Sending join_room request...");
    ws.send(JSON.stringify({ type: 'join_room', payload: { room_id: channelId } }));
}


// --- WebRTC Logic ---

// This function handles publishing our own local media (our UPLINK).
async function startPublishing() {
    try {
        const localUserStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
        localVideoPlayer.value.srcObject = localUserStream;

        publisherPC = new RTCPeerConnection(peerConfiguration);

        localUserStream.getTracks().forEach(track => {
            publisherPC.addTrack(track, localUserStream);
        });

        publisherPC.onicecandidate = (event) => {
            if (event.candidate) {
                ws.send(JSON.stringify({
                    type: 'candidate',
                    payload: { candidate: event.candidate, remote_peer_id: myPeerID }
                }));
            }
        };

        const offer = await publisherPC.createOffer();
        await publisherPC.setLocalDescription(offer);
        console.log("-> Sending publisher offer...");
        ws.send(JSON.stringify({ type: 'offer', payload: { sdp: offer, remote_peer_id: myPeerID } }));

    } catch (err) {
        console.error("Error in startPublishing:", err);
    }
}

// This function handles an offer to subscribe to a remote user's stream (a DOWNLINK).
async function SubscriptionOffer(payload) {
    const remotePeerID = payload.remote_peer_id;
    console.log(`Processing subscription offer for stream from ${remotePeerID}`);

    const subPC = new RTCPeerConnection(peerConfiguration);
    subscriberPCs[remotePeerID] = subPC; // Store the new connection in our map.

    subPC.onicecandidate = event => {
        if (event.candidate) {
            ws.send(JSON.stringify({
                type: 'candidate',
                payload: { candidate: event.candidate, remote_peer_id: remotePeerID }
            }));
        }
    };

    // This event fires when the remote video/audio track arrives.
    subPC.ontrack = (event) => {
        console.log(`Got track from ${remotePeerID}`);
        let remoteVideo = document.getElementById(`video-${remotePeerID}`);
        if (!remoteVideo) {
            remoteVideo = document.createElement('video');
            remoteVideo.id = `video-${remotePeerID}`;
            remoteVideo.autoplay = true;
            remoteVideo.playsInline = true;
            videoGrid.value.appendChild(remoteVideo);
        }
        remoteVideo.srcObject = event.streams[0];
    };

    await subPC.setRemoteDescription(new RTCSessionDescription(payload.sdp));
    const answer = await subPC.createAnswer();
    await subPC.setLocalDescription(answer);
    console.log(`-> Sending answer for subscription to ${remotePeerID}...`);
    ws.send(JSON.stringify({ type: 'answer', payload: { sdp: answer, remote_peer_id: remotePeerID } }));
}

// This function routes incoming ICE candidates to the correct PeerConnection.
async function RemoteCandidate(payload) {
    if (!payload.candidate) return;
    
    const isForPublisher = (payload.remote_peer_id === myPeerID);
    const pc = isForPublisher ? publisherPC : subscriberPCs[payload.remote_peer_id];

    if (pc) {
        await pc.addIceCandidate(new RTCIceCandidate(payload.candidate));
    }
}
</script>

<style>
#video-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
}
video {
    width: 100%;
    height: auto;
    background: black;
    border-radius: 8px;
}
</style>

