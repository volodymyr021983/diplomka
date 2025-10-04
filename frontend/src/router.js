import { createWebHistory, createRouter } from 'vue-router'

import SignUpView from './pages/auth/SignUpView.vue'
import HomeView from './pages/HomeView.vue'
import SignInView from './pages/auth/SignInView.vue'
import TestView from './pages/TestView.vue'
import ResetPwdView from './pages/auth/ResetPwdView.vue'
import SetNewPwdView from './pages/auth/SetNewPwdView.vue'
import EmailVerifyView from './pages/auth/EmailVerifyView.vue'
import CreateServerView from './pages/server/CreateServerView.vue'
import ServerView from './pages/server/ServerView.vue'
import NotFoundView from './pages/NotFoundView.vue'
import AcceptInvitationView from './pages/server/AcceptInvitationView.vue'
import Session from "supertokens-web-js/recipe/session"


async function requiredAuth(to,from,next){
    const hasSession = await Session.doesSessionExist()
    if(hasSession){
        next();
    }else{
        next({name: 'signin'})
    }
}
async function serverAuth(to,form,next){
    const hasSession = await Session.doesSessionExist()
    if(hasSession){
        const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i
        if (!uuidRegex.test(to.params.server_id)) {
        return next({ name: 'notfound' })
      }
      const response = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/api/server/check-connect/${to.params.server_id}/${to.params.channel_id}`)

      if (!response.ok){
        next({name: 'notfound'})
      }

      next();
    }else{
        next({name: 'signin'})
    }
}
const routes = [
    {
        path: '/', 
        name: "home",
        component: HomeView
    },
    {
        path: '/notfound',
        name: "notfound",
        component: NotFoundView
    },
    {
        path: "/auth/signup",
        name: "signup",
        component: SignUpView,
    },
    {
        path: "/auth/signin",
        name: "signin",
        component: SignInView 
    },
    {
        path: "/test",
        name: "test",
        component: TestView,
        beforeEnter: requiredAuth
    },
    {
        path: "/auth/reset-password",
        name: "reset-password",
        component: ResetPwdView
    },
    {
        path: "/auth/set-password",
        name: "set-password",
        component: SetNewPwdView
    },
    {
        path: "/auth/verify-email",
        name: "verify-email",
        component: EmailVerifyView
    },
    {
        path: "/server/create",
        name: "server-create",
        component: CreateServerView,
        beforeEnter: requiredAuth
    },
    {
        path: "/server/:server_id/:channel_id",
        name: "server-connect",
        component: ServerView,
        beforeEnter: serverAuth
    },
    {
        path: "/server/invite/:invite_code",
        name: "invite",
        component: AcceptInvitationView,
        beforeEnter: requiredAuth
    }

    
]

const router = createRouter({history: createWebHistory(),routes,})

export default router