
<template>
    <v-card>
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
                <v-btn icon>
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
        <v-list-item link @click="$emit('ShowProjects')" title="Projects">
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
    </v-navigation-drawer>
    </v-card>
</template>


<script lang="ts">
  import Vue from 'vue'

  interface User {
    ID:          number;
	Username:    string;
    Mail:        string;
    Firstname:   string;
    Lastname:    string;
  }
  const UserInfo: User = {
      ID: 0,
      Username: "Loading",
      Firstname: "Please",
      Lastname: "Wait",
      Mail: "Loading"
  }

  export default Vue.extend({
    name: 'Drawer',
    data: function(){
        return {
            UserInfo,
            ShowUserMenu: false
        }
    },
    mounted: async function(){
        this.UserInfo = await Vue.prototype.$GetRequest("/api/v1/profile")
    }
  })
</script>