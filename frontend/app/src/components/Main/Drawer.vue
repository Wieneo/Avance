
<template>
    <v-navigation-drawer permanent mini-variant app
    >
        <v-list>
        <v-menu :offset-x=true :nudge-width="300" :close-on-content-click="false" v-model="ShowUserMenu" >
            <template v-slot:activator="{ on }">
                    <v-list-item class="px-2">
                    <v-list-item-avatar>
                        <v-img class="ProfilePicture" :src="getProfilePictureLink" v-on="on" @click="ShowUserMenu = true" style="cursor: pointer;" ></v-img>
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
                <v-btn icon @click="GoToSettings">
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

        <v-dialog v-model="ShowAvatarDeleteConfirmation" persistent max-width="380">
            <v-card>
                <v-card-title class="headline">Are you sure?</v-card-title>
                <v-card-text>Do you really want to remove your Profile Picture?<br>This change will take effect immediately!</v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="orange darken-1" text @click="ShowAvatarDeleteConfirmation = false">ABORT</v-btn>
                    <v-btn color="red darken-1" text @click="deleteProfilePicture">DELETE</v-btn>
                </v-card-actions>

            </v-card>
        </v-dialog>
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

  interface JsonError {
      Error:    string;
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
        GoToSettings: function(){
            try{
                this.$router.push({ path: '/settings', query: { setting: "0" } })
            }finally{
                this.ShowUserMenu = false
            }
            
        }
    }
  })
</script>