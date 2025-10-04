import './assets/main.css'
import SuperTokens from 'supertokens-web-js';
import Session from 'supertokens-web-js/recipe/session';
import EmailPassword from 'supertokens-web-js/recipe/emailpassword'
import EmailVerification from "supertokens-web-js/recipe/emailverification";


import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

const backendApi = import.meta.env.VITE_BACKEND_API_URL
const frontendApi = import.meta.env.VITE_FRONTEND_API_URL

SuperTokens.init({
    appInfo: {
        apiDomain: backendApi,
        websiteDomain: frontendApi,
        apiBasePath: "/auth",
        appName: "Discord"
    },
    recipeList:[
        Session.init(),
        EmailPassword.init({signInAndUpFeature:{
            signUpForm:{
                formFields: [{
                    id: "username",
                    label: "Username",
                }]
            }
        }}),
        EmailVerification.init(),
    ],
});

createApp(App).use(router).mount('#app')
