package auth

import (
	"errors"
	"fmt"
	"strings"
	"test/discord/db"
	"time"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func customSignUpPOST(
	originalSignUpPOST func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error), dbContainer *db.DbContainer,
) func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
	return func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
		username := ""
		for _, field := range formFields {
			if field.ID == "username" {
				valueAsString, asStrOk := field.Value.(string)
				if !asStrOk {
					return epmodels.SignUpPOSTResponse{}, errors.New("should never come here as we check the type during validation")
				}
				username = valueAsString
			}
		}
		_, err := getUserUsingUsername(username, dbContainer)
		if err == nil {
			return epmodels.SignUpPOSTResponse{
				GeneralError: &supertokens.GeneralErrorResponse{
					Message: errors.New("username already taken").Error(),
				},
			}, nil
		}

		response, err := originalSignUpPOST(formFields, tenantId, options, userContext)

		if err != nil {
			return epmodels.SignUpPOSTResponse{}, err
		}

		if response.EmailAlreadyExistsError != nil {
			return epmodels.SignUpPOSTResponse{
				GeneralError: &supertokens.GeneralErrorResponse{
					Message: errors.New("email already taken").Error(),
				},
			}, nil

		}

		if response.OK != nil && response.OK.Session != nil {

			err = saveUserUsernameEmail(username, response.OK.User.Email, response.OK.User.ID, dbContainer)
			if err != nil {
				return epmodels.SignUpPOSTResponse{
					GeneralError: &supertokens.GeneralErrorResponse{
						Message: errors.New("unexpected error").Error(),
					},
				}, nil
			}
			response.OK.Session.RevokeSession()

		}
		return response, nil
	}
}

// override functions
func customSendResetEmail(
	originalSendEmail func(input emaildelivery.EmailType, userContext supertokens.UserContext) error, frontend string,
) func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
	return func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		oldResetUrl := fmt.Sprintf("%s/auth/reset-password", frontend)
		newResetUrl := fmt.Sprintf("%s/auth/set-password", frontend)
		input.PasswordReset.PasswordResetLink = strings.Replace(
			input.PasswordReset.PasswordResetLink,
			oldResetUrl,
			newResetUrl, 1,
		)
		return originalSendEmail(input, userContext)
	}
}
func customCreateResetPasswordPOST(
	originalCreateResetPassword func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.GeneratePasswordResetTokenPOSTResponse, error),
	dbContainer *db.DbContainer,
) func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.GeneratePasswordResetTokenPOSTResponse, error) {
	return func(formFields []epmodels.TypeFormField, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.GeneratePasswordResetTokenPOSTResponse, error) {

		ip, err := getIP(options.Req)
		if err != nil {
			return epmodels.GeneratePasswordResetTokenPOSTResponse{}, errors.New("error related to ip address")
		}
		IpAddress, err := getIpAddress(*ip, dbContainer)

		if err != nil {
			errIp := saveIpAddress(*ip, dbContainer)
			if errIp != nil {
				fmt.Println(errIp.Error())
				return epmodels.GeneratePasswordResetTokenPOSTResponse{}, errors.New("error related to ip address")
			}
			IpAddress, _ = getIpAddress(*ip, dbContainer)

		}

		if IpAddress.ResetPwdTime != nil {
			timePassed := time.Now().UnixMilli() - *IpAddress.ResetPwdTime
			if timePassed < 180000 {
				return epmodels.GeneratePasswordResetTokenPOSTResponse{
					GeneralError: &supertokens.GeneralErrorResponse{
						Message: "token can be generated once every 3 minutes",
					},
				}, nil
			}
		}
		response, err := originalCreateResetPassword(formFields, tenantId, options, userContext)
		if response.OK != nil {
			if IpAddress.ResetPwdTime == nil {
				timeNow := time.Now().UnixMilli()
				IpAddress.ResetPwdTime = &timeNow
				dbContainer.DB.Save(&IpAddress)
			} else {
				*IpAddress.ResetPwdTime = time.Now().UnixMilli()
				dbContainer.DB.Save(&IpAddress)
			}
		}
		return response, err
	}
}
func customGenerateVerifyTokenPOST(
	originalGenerateVerifyTokenPOST func(sessionContainer sessmodels.SessionContainer, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.GenerateEmailVerifyTokenPOSTResponse, error), dbContainer *db.DbContainer,
) func(sessionContainer sessmodels.SessionContainer, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.GenerateEmailVerifyTokenPOSTResponse, error) {
	return func(sessionContainer sessmodels.SessionContainer, options evmodels.APIOptions, userContext supertokens.UserContext) (evmodels.GenerateEmailVerifyTokenPOSTResponse, error) {

		ip, err := getIP(options.Req)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, errors.New("error related to ip address")
		}
		IpAddress, err := getIpAddress(*ip, dbContainer)

		if err != nil {
			errIp := saveIpAddress(*ip, dbContainer)
			if errIp != nil {
				fmt.Println(errIp.Error())
				return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, errors.New("error related to ip address")
			}
			IpAddress, _ = getIpAddress(*ip, dbContainer)
		}

		if IpAddress.EmailVerifyTime != nil {
			timePassed := time.Now().UnixMilli() - *IpAddress.EmailVerifyTime
			if timePassed < 180000 {
				return evmodels.GenerateEmailVerifyTokenPOSTResponse{
					GeneralError: &supertokens.GeneralErrorResponse{
						Message: "token can be generated once every 3 minutes",
					},
				}, nil
			}
		}

		response, err := originalGenerateVerifyTokenPOST(sessionContainer, options, userContext)
		if response.OK != nil {
			if IpAddress.EmailVerifyTime == nil {
				timeNow := time.Now().UnixMilli()
				IpAddress.EmailVerifyTime = &timeNow
				dbContainer.DB.Save(&IpAddress)
			} else {
				*IpAddress.EmailVerifyTime = time.Now().UnixMilli()
				dbContainer.DB.Save(&IpAddress)
			}
		}

		return response, err
	}
}
