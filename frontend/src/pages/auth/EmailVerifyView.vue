<template>
<P>Email Verify</P>
<button v-on:click="consumeVerificationCode">Click here to verify your email</button>
</template>

<script setup lang="js">
import { verifyEmail } from "supertokens-web-js/recipe/emailverification";
async function consumeVerificationCode() {
    try {
        const response = await verifyEmail();
        if (response.status === "EMAIL_VERIFICATION_INVALID_TOKEN_ERROR") {
            // This can happen if the verification code is expired or invalid.
            // You should ask the user to retry
            window.alert("Oops! Seems like the verification link expired. Please try again")
        } else {
            // email was verified successfully.
            window.alert("Email Verified succsesful")
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

</style>