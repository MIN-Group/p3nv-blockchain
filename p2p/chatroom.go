package p2p

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/wooyang2018/ppov-blockchain/logger"

	"github.com/wooyang2018/ppov-blockchain/emitter"
)

var topicName = "blockchain"

// ChatRoom represents a subscription to a single PubSub topic. Messages
// can be published to the topic with ChatRoom.Publish, and received
// messages are pushed to the emitter.Emitter.
type ChatRoom struct {
	ctx     context.Context
	ps      *pubsub.PubSub
	topic   *pubsub.Topic
	sub     *pubsub.Subscription
	emitter *emitter.Emitter
	self    peer.ID
}

type ChatMessage []byte

// JoinChatRoom tries to subscribe to the PubSub topic, returning a ChatRoom on success.
func JoinChatRoom(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID) (*ChatRoom, error) {
	// join the pubsub topic
	topic, err := ps.Join(topicName)
	if err != nil {
		return nil, err
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	cr := &ChatRoom{
		ctx:     ctx,
		ps:      ps,
		topic:   topic,
		sub:     sub,
		self:    selfID,
		emitter: emitter.New(),
	}

	// start reading messages from the subscription in a loop
	go cr.readLoop()
	return cr, nil
}

// Publish sends a message to the pubsub topic.
func (cr *ChatRoom) Publish(message ChatMessage) error {
	return cr.topic.Publish(cr.ctx, message)
}

func (cr *ChatRoom) ListPeers() []peer.ID {
	return cr.ps.ListPeers(topicName)
}

// readLoop pulls messages from the pubsub topic and pushes them onto the emitter.Emitter.
func (cr *ChatRoom) readLoop() {
	for {
		msg, err := cr.sub.Next(cr.ctx)
		if err != nil {
			return
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == cr.self {
			continue
		}
		cr.emitter.Emit(msg)
	}
}

const DiscoveryServiceTag = "chatroom"

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h         host.Host
	peerStore *PeerStore
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if !n.peerStore.IsValidID(pi.ID) {
		logger.I().Warnw("found invalid peer", "peerID", pi.ID)
		return
	}
	if err := n.h.Connect(context.Background(), pi); err != nil {
		logger.I().Errorf("failed to connect to peer %s: %v\n", pi.ID, err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host, peerStore *PeerStore) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h, peerStore})
	return s.Start()
}
