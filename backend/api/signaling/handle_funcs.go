package signaling

import "fmt"

func JoinRoom(channel_id string, client *Client) {
	existingChannels.mu.Lock()
	defer existingChannels.mu.Unlock()

	channel := existingChannels.channels[channel_id]

	if channel == nil {
		usersMap := make(map[string]*Client)
		channel = &Channel{
			channel_id: channel_id,
			users:      usersMap,
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
