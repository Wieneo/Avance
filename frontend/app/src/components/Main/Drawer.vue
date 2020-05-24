
<template>
    <v-card>
    <v-navigation-drawer permanent expand-on-hover app
    >
        <v-list>
        <v-list-item class="px-2">
            <v-list-item-avatar>
            <v-img src="https://randomuser.me/api/portraits/women/85.jpg"></v-img>
            </v-list-item-avatar>
        </v-list-item>

        <v-list-item link>
            <v-list-item-content>
            <v-list-item-title class="title">{{UserInfo.Firstname}} {{UserInfo.Lastname}}</v-list-item-title>
            <v-list-item-subtitle>{{UserInfo.Mail}}</v-list-item-subtitle>
            </v-list-item-content>
        </v-list-item>
        </v-list>

        <v-divider>
        </v-divider>

        <v-list
        nav
        dense
        >
        <v-list-item link @click="$emit('ShowProjects')">
            <v-list-item-icon>
            <v-icon>mdi-folder</v-icon>
            </v-list-item-icon>
            <v-list-item-title>Projects</v-list-item-title>
        </v-list-item>
        <v-list-item link>
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
            UserInfo
        }
    },
    mounted: async function(){
        this.UserInfo = await Vue.prototype.$GetRequest("/api/v1/profile")
    }
  })
</script>