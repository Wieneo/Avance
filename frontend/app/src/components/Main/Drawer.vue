
<template>
    <v-navigation-drawer permanent mini-variant app
    >
        <v-list>
        <v-menu :offset-x=true :nudge-width="300" :close-on-content-click="false" v-model="ShowUserMenu" >
            <template v-slot:activator="{ on }">
                    <v-list-item class="px-2">
                    <v-list-item-avatar>
                        <v-img class="ProfilePicture" src="https://randomuser.me/api/portraits/women/85.jpg" v-on="on" @click="ShowUserMenu = true" style="cursor: pointer;" ></v-img>
                    </v-list-item-avatar>
                </v-list-item>
            </template>
            <v-container style="background-color: white;">
            <v-row>
                <v-col style="min-width: 300px; max-width: 300px;">
                    <v-card style="text-align: center; overflow-y: hidden;" flat>
                        <v-card-title primary-title class="justify-center">{{UserInfo.Firstname}} {{UserInfo.Lastname}}</v-card-title>
                        <v-card-subtitle>{{UserInfo.Username}}</v-card-subtitle>
                        <v-card-text class="TicketDisplayProperty"><v-btn class="ma-2" outlined color="green"><v-icon style="margin-right: 10px;">mdi-circle</v-icon> Online</v-btn></v-card-text>
                    </v-card>
                </v-col>
                <v-btn icon @click="ShowEditMenu = true; ShowUserMenu = false">
                    <v-icon>mdi-cog</v-icon>
                </v-btn>
            </v-row>
            </v-container>
        </v-menu>
        </v-list>

        <v-divider>
        </v-divider>

        <v-list
        nav
        dense
        >
        <v-list-item link @click="showProjects = true" title="Projects">
            <v-list-item-icon>
            <v-icon>mdi-folder</v-icon>
            </v-list-item-icon>
            <v-list-item-title>Projects</v-list-item-title>
        </v-list-item>
        <v-list-item link title="Starred">
            <v-list-item-icon>
            <v-icon>mdi-star</v-icon>
            </v-list-item-icon>
            <v-list-item-title>Starred</v-list-item-title>
        </v-list-item>
        </v-list>


        <v-dialog v-model="ShowEditMenu" persistent max-width="350">
            <v-card>
                <v-card-title class="headline">Your Profile</v-card-title>
                <v-card-text>
                    <div style="text-align: center;">
                        <v-list-item-avatar style="width: 128px; height: 128px;">
                            <v-img class="ProfilePicture" src="https://randomuser.me/api/portraits/women/85.jpg" style="cursor: pointer; " ></v-img>
                        </v-list-item-avatar>
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
                <v-btn color="orange darken-1" text @click="ShowEditMenu = false">Abort</v-btn>
                <v-btn color="green darken-1" text @click="SaveChanges" :loading="Loading">Save</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
        <template v-slot:append>
            <v-list>
                <v-list-item link title="Settings" @click="GoToSettings">
                    <v-list-item-icon>
                        <v-icon>mdi-cog</v-icon>
                    </v-list-item-icon>
                    <v-list-item-title>Settings</v-list-item-title>
                </v-list-item>
            </v-list>
        </template>
        <ProjectsContainer v-bind:showProjects="showProjects" v-on:closeProjects="showProjects = false"/>
    </v-navigation-drawer>
</template>


<script lang="ts">
  import Vue from 'vue'
  import ProjectsContainer from '../misc/ProjectsContainer.vue';

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
      Username: "Loading",
      Firstname: "Please",
      Lastname: "Wait",
      Mail: "Loading",
      Password: ""
  }

  //Initialize seperately so we don't create a reference
  const ChangedProfileInfo: User = {
      ID: 0,
      Username: "Loading",
      Firstname: "Please",
      Lastname: "Wait",
      Mail: "Loading",
      Password: ""
  }

  export default Vue.extend({
    name: 'Drawer',
    components:{
        ProjectsContainer
    },

    data: function(){
        return {
            showProjects: false,
            Loading: false,
            UserInfo,
            ChangedProfileInfo,
            NewPassword1: "",
            NewPassword2: "",
            ShowUserMenu: false,
            ShowEditMenu: false,
            rules: [
                value => !!value || 'Required.',
            ] as ((value: any) => true | "Required.")[],
        }
    },
    watch:{
        ShowEditMenu (to){
            if (!to){
                //CleanUP Form
                console.log("Resetting Inputs from Profile Edit")
                Object.assign(this.ChangedProfileInfo, this.UserInfo)
                this.NewPassword1 = ""
                this.NewPassword2 = ""
            }
        }
    },
    mounted: async function(){
        const data =  await Vue.prototype.$GetRequest("/api/v1/profile")
        //Assign so we dont create a reference here
        Object.assign(this.UserInfo, data);
        Object.assign(this.ChangedProfileInfo, data);
    },
    methods:{
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
                this.ShowEditMenu = false
            }
            this.Loading = false
        },
        GoToSettings: function(){
            this.$router.push({ path: '/settings', query: {} })
        }
    }
  })
</script>