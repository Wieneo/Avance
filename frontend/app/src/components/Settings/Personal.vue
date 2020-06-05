<template>
    <v-container>
        <v-row>
            <v-col>
                <v-card>
                    <v-card-title class="headline">Your Profile</v-card-title>
                    <v-card-text>
                        <div style="text-align: center;">
                            <v-list-item-avatar style="width: 128px; height: 128px;">
                                <v-img class="ProfilePicture" :src="getProfilePictureLink" style="cursor: pointer; " >
                                    <div class="EditProfilePicture" @click="uploadProfilePicture" >
                                        EDIT
                                        <form action="/api/v1/profile/avatar" method="POST" target="submitDeflector" enctype="multipart/form-data">
                                            <input type="file" style="display: none;" onchange="this.form.submit()" name="avatar" id="uploader">
                                        </form>
                                    </div>
                                </v-img>
                            </v-list-item-avatar>
                            <br>
                            <v-btn title="delete profile picture" color="orange darken-1" icon @click="ShowAvatarDeleteConfirmation = true;"><v-icon>mdi-delete-outline</v-icon></v-btn>
                        </div>
                        <v-text-field label="Username" :rules="rules" hide-details="auto" v-model="ChangedProfileInfo.Username"></v-text-field>
                        <v-text-field style="margin-top: 20px;" label="Firstname" :rules="rules" hide-details="auto" v-model="ChangedProfileInfo.Firstname"></v-text-field>
                        <v-text-field style="margin-top: 20px;" label="Lastname" :rules="rules" hide-details="auto" v-model="ChangedProfileInfo.Lastname"></v-text-field>
                        <v-text-field style="margin-top: 20px;" label="E-Mail" :rules="rules" hide-details="auto" v-model="ChangedProfileInfo.Mail"></v-text-field>
                        <v-text-field style="margin-top: 25px;" label="Password" hide-details="auto" type="password" v-model="NewPassword1"></v-text-field>
                        <span>Only needed when changing password!</span>
                        <v-text-field label="Repeat Password" hide-details="auto" v-model="NewPassword2" type="password"></v-text-field>
                    </v-card-text>
                    <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="orange darken-1" text @click="ResetInput">Reset</v-btn>
                    <v-btn color="green darken-1" text @click="SaveChanges" :loading="Loading">Save</v-btn>
                    </v-card-actions>
                </v-card>
            </v-col>
        </v-row>
        <iframe name="submitDeflector" style="display: none;" id="submitDeflector" @load="HandleUploadFinish"> </iframe>
    </v-container>
</template>
<script lang="ts">
import Vue from 'vue'

interface JsonError {
    Error:    string;
}

interface User {
ID:          number;
Username:    string;
Mail:        string;
Firstname:   string;
Lastname:    string;
Password:    string;
}

const UserInfo: User = {
    ID: 0,
    Username: "",
    Firstname: "",
    Lastname: "",
    Mail: "",
    Password: ""
}

//Initialize seperately so we don't create a reference
const ChangedProfileInfo: User = {
    ID: 0,
    Username: "",
    Firstname: "",
    Lastname: "",
    Mail: "",
    Password: ""
}

export default Vue.extend({
    name: 'Personal',
    data: function(){
        return {
            NewPassword1: "",
            NewPassword2: "",
            Loading: false,
            UserInfo,
            ChangedProfileInfo,
            ShowAvatarDeleteConfirmation: false,
            rules: [
                value => !!value || 'Required.',
            ] as ((value: any) => true | "Required.")[],
            IFrameLoaded: false,
        }
    },
    mounted: async function(){
        const data =  await Vue.prototype.$Request("GET", "/api/v1/profile")
        //Assign so we dont create a reference here
        Object.assign(this.UserInfo, data);
        Object.assign(this.ChangedProfileInfo, data);
    },
    computed: {
        // a computed getter
        getProfilePictureLink: function () {
            return '/api/v1/profile/avatar?' + performance.now()
        }
     },
    methods:{
        ResetInput: function(){
            Object.assign(this.ChangedProfileInfo, this.UserInfo)
            this.NewPassword1 = ""
            this.NewPassword2 = ""
        },
        SaveChanges: async function(){
            if (this.NewPassword1.trim().length > 0){
                if (this.NewPassword1 == this.NewPassword2){
                    this.ChangedProfileInfo.Password = this.NewPassword1
                }else{
                    Vue.prototype.$NotifyError("Passwords don't match")
                    return
                }
            }

            this.Loading = true
            const newinfo = await Vue.prototype.$Request("PATCH", "/api/v1/profile", this.ChangedProfileInfo)
            if (newinfo.Error == undefined){
                Object.assign(this.UserInfo, newinfo)
                Vue.prototype.$NotifySuccess("Profile updated")
            }
            this.Loading = false
        },
        deleteProfilePicture: async function(){
            await Vue.prototype.$Request("DELETE", "/api/v1/profile/avatar", null)
            const imgElement: HTMLCollectionOf<HTMLElement> = document.getElementsByClassName('v-image__image') as HTMLCollectionOf<HTMLElement>
            for (let i = 0 ; i < imgElement.length; i++){
                if (imgElement[i].style.backgroundImage.includes("avatar")) {
                    imgElement[i].style.backgroundImage = (imgElement[i].style.backgroundImage.substr(0, imgElement[i].style.backgroundImage.length - 2)) + "?" + (new Date()) + `")`
                }
            }
            this.ShowAvatarDeleteConfirmation = false
        },

        uploadProfilePicture: async function(){
            const element: HTMLElement = document.getElementById("uploader") as HTMLElement
            element.click()
        },

        HandleUploadFinish: async function(){
            if (this.IFrameLoaded == false) {
                this.IFrameLoaded = true
                return
            }
            const element: HTMLIFrameElement = document.getElementById("submitDeflector") as HTMLIFrameElement
            const elementDocument: Document = element.contentDocument as Document
            if (elementDocument.documentElement.innerText.length > 0) {
                const JSONData: JsonError = JSON.parse(elementDocument.documentElement.innerText)
                Vue.prototype.$NotifyError(JSONData.Error)
            } else {
                Vue.prototype.$NotifySuccess("Avatar Updated")
                const imgElement: HTMLCollectionOf<HTMLElement> = document.getElementsByClassName('v-image__image') as HTMLCollectionOf<HTMLElement>
                for (let i = 0 ; i < imgElement.length; i++){
                    if (imgElement[i].style.backgroundImage.includes("avatar")) {
                        if (imgElement[i].style.backgroundImage.includes("?")){
                            imgElement[i].style.backgroundImage = imgElement[i].style.backgroundImage.split("?")[0] + `")`
                        }
                        imgElement[i].style.backgroundImage = (imgElement[i].style.backgroundImage.substr(0, imgElement[i].style.backgroundImage.length - 2)) + "?" + performance.now() + `")`
                    }
                }
            }
        }
    }
})
</script>