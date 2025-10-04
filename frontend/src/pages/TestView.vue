<template>
    <button @click="doesSessionExist">
    Check</button><br>
    <input v-model="userId">
    <input v-model="isVerified">
    <button type="button" @click="Logout">Sign Out</button>
    <button type="button" @click="sendEmail">Send Verification Link</button>
    <button type="button" @click="GetEmailVerification">Check if email verified</button>

</template>


<script setup lang="js">
import Session from 'supertokens-web-js/recipe/session'
import { sendVerificationEmail,isEmailVerified } from "supertokens-web-js/recipe/emailverification";

import { ref } from 'vue'
const userId = ref('')
const jwt = ref('')
const isVerified = ref(false)

async function Logout(){
    await Session.signOut();
    window.location.href="/"
}
async function GetEmailVerification(){
    const response = await isEmailVerified()
    isVerified.value = response.isVerified
}
async function doesSessionExist(){
    if (await Session.doesSessionExist()){
        userId.value = await Session.getUserId()
        jwt.value = await Session.getAccessToken()
    }else{
        console.log("NOT EXISTS")
    }
}
async function sendEmail() {
    try {
        const response = await sendVerificationEmail();
        if (response.status === "EMAIL_ALREADY_VERIFIED_ERROR") {
            // This can happen if the info about email verification in the session was outdated.
            // Redirect the user to the home page
            window.location.assign("/home");
        } else {
            // email was sent successfully.
            window.alert("Please check your email and click the link in it")
        }
    } catch (err) {
        if (err.isSuperTokensGeneralError === true) {
            // this may be a custom error message sent from the API by you.
            window.alert(err.message);
        } else {
            window.alert("Oops! Something went wrong.");
        }
    }
}
</script>
<style scoped>
.exist{
    height: 600px;
    width: 200px;
    background-color: aqua;
}
.nope{
    height: 100px;
    width: 100px;
    background-color: brown;
}
</style>

