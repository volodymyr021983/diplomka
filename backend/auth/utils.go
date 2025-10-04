package auth

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"test/discord/db"
	"test/discord/db/models"
)

func getUserUsingUsername(username string, dbContainer *db.DbContainer) (*models.UserProfile, error) {
	var UserProfile models.UserProfile
	result := dbContainer.DB.Limit(1).Find(&UserProfile, "username = ?", username)
	if result.RowsAffected != 1 {
		return nil, errors.New("username not found")
	}
	return &UserProfile, nil
}
func saveUserUsernameEmail(username string, email string, userID string, dbContainer *db.DbContainer) error {
	fmt.Print("In Save Users")
	userProfile := models.UserProfile{
		UserID:   userID,
		Username: username,
		Email:    email,
	}
	result := dbContainer.DB.Create(&userProfile)

	if result.Error != nil {
		fmt.Println("Error occur while adding username to db")
		return result.Error
	}

	if result.RowsAffected != 1 {
		fmt.Println("Error added more then one user")
		return errors.New("added more then one user")
	}

	return nil
}

func getUserUsingEmail(email string, dbContainer *db.DbContainer) (*models.UserProfile, error) {
	var UserProfile models.UserProfile
	result := dbContainer.DB.Limit(1).Find(&UserProfile, "email = ?", email)

	if result.RowsAffected != 1 {
		return nil, errors.New("email not found")
	}
	return &UserProfile, nil
}
func getIpAddress(ip string, dbContainer *db.DbContainer) (*models.IpAddress, error) {
	var IpAddress models.IpAddress
	result := dbContainer.DB.Limit(1).Find(&IpAddress, "ip_address = ?", ip)

	if result.RowsAffected != 1 {
		return nil, errors.New("ip_address not found")
	}
	return &IpAddress, nil

}
func saveIpAddress(ip string, dbContainer *db.DbContainer) error {
	IpAddress := models.IpAddress{
		IpAddress: ip,
	}
	result := dbContainer.DB.Create(&IpAddress)

	if result.Error != nil {
		fmt.Println("Error occur while adding ipAddress to db")
		return result.Error
	}

	if result.RowsAffected != 1 {
		fmt.Println("Error added more then one ipAddr")
		return errors.New("added more then one ipAddr")
	}
	return nil
}
func getIP(req *http.Request) (*string, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return &ip, nil
}
