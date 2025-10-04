<template>
    <div class="form-wrapper centered">
        <form class="centered" @submit.prevent="SetNewPwdClicked">
            <input v-model="password" placeholder="Password"required type="Password"><br>
            <button type="submit">Sign Up</button>
        </form>
    </div>
</template>

<script setup lang="js">
import { submitNewPassword } from "supertokens-web-js/recipe/emailpassword";
import { ref } from 'vue'
const statusik = ref('')
const password = ref('')
async function SetNewPwdClicked(){
    try{
        const response = await submitNewPassword({
            formFields:
            [{
                id: "password",
                value: password.value
            }]
        })
        statusik.value = response.status
        if (response.status === "FIELD_ERROR"){
            response.formFields.forEach(formField => {
            console.log("fields errors")
            console.log(formField.error)
            if (formField.id === "password"){
                console.log("EMAIL ERROR")
                window.alert(formatField.error)
                }
            })
        }else if(response.status === "RESET_PASSWORD_INVALID_TOKEN_ERROR"){
            window.alert(response.reason)
            console.log("RESET_PASSWORD_INVALID_TOKEN_ERROR")
        }else{
           console.log("ddd")
           console.log(statusik)
           console.log(response.reason)
        }
    }
        catch (err){
        if(err.isSuperTokensGeneralError === true){
            window.alert(err.message);
        }else{
            console.log(err.message)
            window.alert("Oops! Something went wrong.");
        }
    }

}
</script>