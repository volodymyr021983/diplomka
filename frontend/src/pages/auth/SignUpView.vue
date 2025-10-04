<template>
    <div class="form">
        <form @submit.prevent="SignUpClicked">
            <fieldset class="fieldset">
            <legend>Email:</legend>
            <input class="input" v-model="email" placeholder="Email" type="Email" required><br>
            </fieldset>
            <p v-if="emailError">{{ emailError }}</p>
           
            
            <fieldset class="fieldset">
            <legend>Username:</legend>
            <input class="input" v-model="username" placeholder="Username" type="Text" required><br>
            </fieldset>
            <p v-if="usernameError">{{ usernameError }}</p>
            

            <fieldset class="fieldset">
            <legend>Password:</legend>
            <input class="input" v-model="password" placeholder="Password" type="Password" required><br>
            </fieldset>
            <p v-if="passwordError">{{ passwordError }}</p>
            

            <button type="submit"> Sign Up </button>
        
        </form>
    </div>
</template>
<script setup lang="js">
import { signUp } from "supertokens-web-js/recipe/emailpassword";
import { ref } from 'vue'

const email = ref('')
const password = ref('')
const username = ref('')

const emailError = ref('')
const usernameError = ref('')
const passwordError = ref('')
const generalError = ref('')

async function SignUpClicked(){
        emailError.value = ''
        usernameError.value = ''
        passwordError.value = ''    
        generalError.value = ''
   try{
        const response = await signUp(
            {
            formFields:
            [{
                id: "email",
                value: email.value
            },
            {
                id: "password",
                value: password.value
            },
            {
                id: "username",
                value: username.value
            }]
        })
        if (response.status === "FIELD_ERROR"){
        response.formFields.forEach(formField => {
            if (formField.id === "email"){
                emailError.value = "Email exists. Please sign in instead"
            }else if(formField.id === "username"){
                usernameError.value = "Username has been already taken"
            }else if( formField.id === "password"){
                passwordError.value = formField.error
            }
        })
    }else if(response.status === "SIGN_UP_NOT_ALLOWED"){
        generalError.value = "Unexpected error"
    }

}catch(err){
        if(err.isSuperTokensGeneralError === true){
            if(err.message === "username already taken"){
                usernameError.value = "Username has been already taken"
            }
            if(err.message === "email already taken"){
                emailError.value = "Email exists. Please sign in instead"
            }
            if(err.message === "unexpected error"){
                generalError.value = "Unexpected error"
            }
        }else{
            generalError.value = "Unexpected error"
        }
    }
}
    
</script>

<style scoped>
.fieldset {
  --main-color: #4a4d96;
  border: 2px var(--main-color) solid;
  font-size: 0.75rem;
  border-radius: 5px;
}
legend {
  color: var(--main-color);
}
.fieldset .input{
  outline: none;
  border: none;
  background: transparent;

  padding: 0 10px 10px 10px;
  font-size: 0.75rem;
}
form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
button {
 font-size: 17px;
 padding: 0.5em 2em;
 margin-top: 5px;
 border: transparent;
 box-shadow: 2px 2px 4px rgba(0,0,0,0.4);
 background: black;
 color: white;
 border-radius: 4px;
 transition: background-color 0.3s ease; /* smooth change */
}

button:hover {
 background-color: rgb(54, 55, 55);
}

button:active {
 transform: translate(0em, 0.2em);
}
</style>