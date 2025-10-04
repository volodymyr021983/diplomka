<template>
    <div class="form-wrapper centered">
        <form class="centered" @submit.prevent="ResetClicked">
            <input v-model="email" placeholder="Email" required type="Email"><br>
            <button type="submit">Reset</button>
        </form>
    </div>
</template>

<script setup lang="js">
import { sendPasswordResetEmail } from "supertokens-web-js/recipe/emailpassword";
import { ref } from 'vue'

const email = ref('')


async function ResetClicked(){
    try{
        const response = await sendPasswordResetEmail({
            formFields:[{
                id: "email",
                value: email.value
            }]

        })
        console.log("Send password called")
        if (response.status === "FIELD_ERROR"){
            response.formFields.forEach(formField => {
            console.log("fields errors")
            console.log(formField.error)
            if (formField.id === "email"){
                console.log("EMAIL ERROR")
                window.alert(formatField.error)
                }
            })
        }else if(response.status === "RESET_PASSWORD_INVALID_TOKEN_ERROR"){
            window.alert(response.reason)
        }else{
            window.alert("GOOOD!")
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

<style scoped>
.centered{
    margin: 0;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
}
.form-wrapper{
    background-color: white;
    width: 600px;
    height: 600px;
}
input{
    width: 200px;
    height: 50px;
    margin: 15px;
}
</style>