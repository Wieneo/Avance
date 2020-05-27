<template>
    <div>
        <v-row
          align="center"
          justify="center"
        >
          <v-col
            cols="12"
            sm="8"
            md="4"
          >
            <v-card class="elevation-12">
              <v-toolbar
                color="primary"
                dark
                flat
              >
                <v-toolbar-title>Hostname</v-toolbar-title>
                <v-spacer />
              </v-toolbar>
              <v-card-text>
                <v-form>
                  <v-text-field
                    label="Username"
                    name="login"
                    prepend-icon="mdi-account"
                    type="text"
                    v-model="Username"
                    @keyup.enter.native="Login"
                    :error="CredentialsInvalid"
                  />

                  <v-text-field
                    id="password"
                    label="Password"
                    name="password"
                    prepend-icon="mdi-form-textbox-password"
                    type="password"
                    v-model="Password"
                    @keyup.enter.native="Login"
                    :error="CredentialsInvalid"
                  />
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-spacer />
                <v-btn color="primary" @click="Login" :loading=Loading>Login</v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
    </div>
</template>


<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
    name: "Login",
    mounted: async function(){
        //Check if session already is valid -> Redirect to main page
        if ((await Vue.prototype.$GetRequest("/api/v1/session")).Authorized){
            //window.location.href = "/"
        }
    },
    data: function(){
        return {
          Username: "",
          Password: "",
          Loading: false,
          CredentialsInvalid: false
        }
    },
    methods:{
      Login: async function (){
        if (this.Username.length > 0 && this.Password.length > 0){
          this.CredentialsInvalid = false
          this.Loading = true
          const Result = await Vue.prototype.$Request("POST", "/api/v1/login", {Username: this.Username, Password: this.Password}, true)
          if (Result.Error == undefined){
            Vue.prototype.$SetCookie('session', Result.SessionKey, 3600)
            window.location.href = "/"
          }else{
            this.CredentialsInvalid = true
          }
        }
        this.Loading = false
      }
    }
})
</script>