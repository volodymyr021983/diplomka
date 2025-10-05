package channels

import (
	"errors"
	"fmt"
	"test/discord/db"
	"test/discord/db/models"

	"github.com/google/uuid"
	"github.com/olahol/melody"
)

func CreateTextChannel(channelModel models.Channels, dbContainer *db.DbContainer) error {
	result := dbContainer.DB.Create(&channelModel)

	if result.Error != nil {
		return errors.New("unexpected error occurs")
	}
	if result.RowsAffected != 1 {
		return errors.New("unexpected error occurs")
	}
	return nil
}

func FindChannelById(channel_id string, dbContainer *db.DbContainer) *models.Channels {
	var channel models.Channels
	result := dbContainer.DB.Limit(1).Find(&channel, "channel_id = ?", channel_id)

	if result.RowsAffected != 1 {
		return nil
	}
	return &channel
}

func GetNewChannelId(dbContainer *db.DbContainer) (*string, error) {
	channel_id := uuid.New().String()

	result := FindChannelById(channel_id, dbContainer)

	if result != nil {
		return nil, errors.New("server already exists")
	}
	return &channel_id, nil
}

func GetFirstChannel(server_id string, dbContainer *db.DbContainer) (*string, error) {
	var channel models.Channels
	result := dbContainer.DB.Limit(1).Find(&channel, "own_server_id = ?", server_id)

	if result.RowsAffected != 1 {
		return nil, errors.New("cant find channel")
	}
	return &channel.ChannelId, nil
}
func GetServerChannels(server_id string, dbContainer *db.DbContainer) []channelsResponse {
	var channels []channelsResponse

	result := dbContainer.DB.Model(&models.Channels{}).Select("Channelname", "ChannelType", "ChannelId").Where("own_server_id = ?",
		server_id).Find(&channels)

	if result.RowsAffected == 0 {
		fmt.Println("nil")
		return nil
	}
	return channels
}

func getUserUsernameUsingId(user_id string, dbContainer *db.DbContainer) *string {
	var UserProfile models.UserProfile
	result := dbContainer.DB.Limit(1).Find(&UserProfile, "user_id = ?", user_id)
	if result.RowsAffected != 1 {
		return nil
	}
	return &UserProfile.Username
}
func deleteChannel(channel models.Channels, m *melody.Melody, dbContainer *db.DbContainer) error {
	serverChannels := getChannelCount(channel.OwnServerId, dbContainer)
	if serverChannels[channel.ChannelType] <= 1 {
		return errors.New("must be at least 1 channel")
	}
	WSDisconnectAllChannelSession(m, channel.ChannelId)
	result := dbContainer.DB.Delete(&channel)
	if result.RowsAffected != 1 {
		return errors.New("unexpected error")
	}
	return nil
}
func getChannelCount(server_id string, dbContainer *db.DbContainer) map[string]int {
	var channels []models.Channels
	dbContainer.DB.Find(&channels, "own_server_id = ?", server_id)
	var textChannelsCount int
	var voiceChannelsCount int

	for _, channel := range channels {
		if channel.ChannelType == "text" {
			textChannelsCount++
		} else {
			voiceChannelsCount++
		}
	}
	return map[string]int{"text": textChannelsCount, "voice": voiceChannelsCount}
}
func WSDisconnectAllChannelSession(m *melody.Melody, channel_id string) error {
	sessions, err := m.Sessions()

	if err != nil {
		return errors.New("unexpected error")
	}
	for _, session := range sessions {
		ses_channel, _ := session.Get("channel_id")
		if ses_channel == channel_id {
			session.Close()
		}
	}
	return nil
}
func WSDisconnectFromAllChannels(m *melody.Melody, server_id string, dbContainer *db.DbContainer) error {
	channels := GetServerChannels(server_id, dbContainer)
	for _, channel := range channels {
		err := WSDisconnectAllChannelSession(m, channel.ChannelId)
		if err != nil {
			return err
		}
	}
	return nil
}

type channelsResponse struct {
	Channelname string
	ChannelId   string
	ChannelType string
}
