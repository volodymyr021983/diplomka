package servers

import (
	"errors"
	"os"
	"test/discord/db"
	"test/discord/db/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func FindServerById(server_id string, dbContainer *db.DbContainer) *models.Servers {
	var server models.Servers
	result := dbContainer.DB.Limit(1).Find(&server, "server_id = ?", server_id)

	if result.RowsAffected != 1 {
		return nil
	}
	return &server
}

func GetNewServerId(dbContainer *db.DbContainer) (*string, error) {
	server_id := uuid.New().String()

	result := FindServerById(server_id, dbContainer)

	if result != nil {
		return nil, errors.New("server already exists")
	}
	return &server_id, nil
}

func CreateNewServer(serverModel *models.Servers, dbContainer *db.DbContainer) error {
	result := dbContainer.DB.Create(&serverModel)

	if result.Error != nil {
		return errors.New("unexpected error occurs")
	}
	if result.RowsAffected != 1 {
		return errors.New("unexpected error occurs")
	}
	return nil
}

func GetUserServers(user_id string, dbContainer *db.DbContainer) ([]GetServersResponse, error) {
	var servers []GetServersResponse
	err := dbContainer.DB.Model(models.Servers{}).
		Joins("JOIN server_members ON server_members.server_id = servers.server_id").
		Where("server_members.user_id = ?", user_id).
		Find(&servers).Error
	if err != nil {
		return nil, errors.New("unexpected error")
	}
	if len(servers) == 0 {
		return nil, errors.New("not found")
	}
	return servers, nil
}

func IsMember(user_id string, server_id string, dbContainer *db.DbContainer) bool {
	var server_member models.ServerMembers

	result := dbContainer.DB.Where("server_id = ? AND user_id = ?", server_id, user_id).Find(&server_member)
	return result.RowsAffected == 1
}

func CreateInviteCode(server_id string) (string, *int64, *int64, error) {
	envErr := godotenv.Load()
	if envErr != nil {
		return "", nil, nil, errors.New("error loading env")
	}

	secret_key := []byte(os.Getenv("SECRET_KEY"))
	create_at := time.Now().Unix()
	exp := time.Now().Add(time.Minute * 15).Unix()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": server_id,
		"exp": exp,
	})

	tokenString, err := claims.SignedString(secret_key)

	if err != nil {
		return "", nil, nil, err
	}
	return tokenString, &create_at, &exp, nil
}
func SaveInviteToken(token models.InvitationCodes, dbContainer *db.DbContainer) error {
	result := dbContainer.DB.Create(&token)

	if result.Error != nil {
		return errors.New("unexpected error occurs")
	}
	if result.RowsAffected != 1 {
		return errors.New("unexpected error occurs")
	}
	return nil
}
func IsServerInviteCodeExists(server_id string, dbContainer *db.DbContainer) (bool, models.InvitationCodes) {
	var invite_code models.InvitationCodes

	result := dbContainer.DB.Where("server_id = ?", server_id).Find(&invite_code)
	return result.RowsAffected == 1, invite_code
}
func DeleteInviteCode(invite_code models.InvitationCodes, dbContainer *db.DbContainer) error {
	result := dbContainer.DB.Delete(&invite_code)
	if result.RowsAffected != 1 {
		return errors.New("error during invite code deletion")
	}
	return nil
}
func VerifyInviteCode(invite_code string, dbContainer *db.DbContainer) (*jwt.Token, error) {

	token, err := jwt.Parse(invite_code, func(token *jwt.Token) (interface{}, error) {
		envErr := godotenv.Load()
		if envErr != nil {
			return nil, nil
		}
		secret_key := []byte(os.Getenv("SECRET_KEY"))
		return secret_key, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, errors.New("error during jwt parsing")
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	return token, nil

}
func AddNewUser(server_member models.ServerMembers, dbContainer *db.DbContainer) error {
	result := dbContainer.DB.Create(&server_member)

	if result.Error != nil {
		return errors.New("unexpected error occurs")
	}
	if result.RowsAffected != 1 {
		return errors.New("unexpected error occurs")
	}
	return nil
}

type APIServerResponse struct {
	OK  bool    `json:"OK"`
	Err *string `json:"error"`
}
type ServerNameBody struct {
	Servername string
}
type GetServersResponse struct {
	Servername string
	ServerId   string
}
