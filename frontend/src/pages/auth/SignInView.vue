<template>
    <div class="form-wrapper centered">
        <form class="centered" @submit.prevent="SignInClicked">
            <input v-model="email" placeholder="Email" required type="email"><br>
            <input v-model="password" placeholder="Password"required type="Password"><br>
            <button type="submit">Sign In</button>
        </form>
    </div>
</template>
<script setup lang="js">
import { signIn } from "supertokens-web-js/recipe/emailpassword";
import { ref } from 'vue'

const email = ref('')
const password = ref('')

const credentialsError = ref('')
const generalError = ref('')
async function SignInClicked(){
    try{
        const response = await signIn({
            formFields: [{
                id: "email",
                value: email.value
            },{
                id: "password",
                value: password.value
            }]
        })
        if(response.status === "FIELD_ERROR"){
            console.log(form)
            response.formFields.forEach(formField => {
                console.log(formField.error)
                if(formField.id === "email"){
                    window.alert(formField.error)
                }
                if(formField.id === "password"){
                    window.alert(formField.error)
                }
            })
        }else if(response.status === "WRONG_CREDENTIALS_ERROR"){
            credentialsError.value = "Email password combination is incorrect."
        }else if(response.status === "SIGN_IN_NOT_ALLOWED"){
            generalError.value = "Unexpected error occurs"
        }
    }catch(err){
        if(err.isSuperTokensGeneralError === true){
            window.alert(err.message);
        }else{
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