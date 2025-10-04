package auth

import (
	"log"
	"os"

	"test/discord/db"

	"github.com/joho/godotenv"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Supertokens_init(DbContainer *db.DbContainer) {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	connectUri := os.Getenv("CONNECTION_URI")
	backendUrl := os.Getenv("BACKEND_API_URL")
	frontend := os.Getenv("FRONTEND_API_URL")
	appName := os.Getenv("APP_NAME")

	smtpSettings := smtp_init()

	apiBasePath := "/auth"
	websiteBasePath := "/auth"

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectUri,
		},
		AppInfo: supertokens.AppInfo{
			AppName:         appName,
			APIDomain:       backendUrl,
			WebsiteDomain:   frontend,
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						//SignUp override
						originalSignUpPOST := *originalImplementation.SignUpPOST

						(*originalImplementation.SignUpPOST) = customSignUpPOST(originalSignUpPOST, DbContainer)
						//Reset Token override
						originalCreateResetPasswordPost := *originalImplementation.GeneratePasswordResetTokenPOST

						(*originalImplementation.GeneratePasswordResetTokenPOST) = customCreateResetPasswordPOST(originalCreateResetPasswordPost, DbContainer)

						return originalImplementation
					},
				},
				EmailDelivery: &emaildelivery.TypeInput{
					Service: emailpassword.MakeSMTPService(emaildelivery.SMTPServiceConfig{
						Settings: smtpSettings,
					}),
					Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {

						originalSendEmail := *originalImplementation.SendEmail

						(*originalImplementation.SendEmail) = customSendResetEmail(originalSendEmail, frontend)

						return originalImplementation
					},
				},
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID:       "username",
							Validate: usernameValidator,
						},
					},
				},
			}),
			session.Init(nil),
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeRequired,
				EmailDelivery: &emaildelivery.TypeInput{
					Service: emailverification.MakeSMTPService(emaildelivery.SMTPServiceConfig{
						Settings: smtpSettings,
					}),
				},
				Override: &evmodels.OverrideStruct{
					APIs: func(originalImplementation evmodels.APIInterface) evmodels.APIInterface {
						originalGenerateVerifyTokenPOST := *originalImplementation.GenerateEmailVerifyTokenPOST

						(*originalImplementation.GenerateEmailVerifyTokenPOST) = customGenerateVerifyTokenPOST(originalGenerateVerifyTokenPOST, DbContainer)

						return originalImplementation
					},
				},
			}),
		},
	})

	if err != nil {
		panic(err.Error())
	}
}
