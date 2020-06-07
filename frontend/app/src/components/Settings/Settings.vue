<template>
  <div style="height: 100%; width: 94%">
    <Drawer v-on:ShowProjects="showProjects = true" />
        <v-tabs background-color="#46494c" dark style="margin-left: 56px; height: 100%;" vertical v-model="SelectedTab" @change="UpdateQuery">
          <v-tab>
            <v-icon left>mdi-account</v-icon>Personal
          </v-tab>
          <v-tab>
            <v-icon left>mdi-bell-ring</v-icon>Notifications
          </v-tab>

          <v-tab v-if="Permissions.Admin">
            <v-icon left>mdi-cog</v-icon>Instance
          </v-tab>

          <v-tab v-if="Permissions.Admin">
            <v-icon left>mdi-application</v-icon>Deployment
          </v-tab>

          <v-tab v-if="Permissions.Admin">
            <v-icon left>mdi-human</v-icon>User & Groups
          </v-tab>

          <v-tab v-if="Permissions.Admin">
            <v-icon left>mdi-ticket</v-icon>Ticket
          </v-tab>

          <v-tab v-if="Permissions.Admin">
            <v-icon left>mdi-transit-connection-variant</v-icon>Integrations
          </v-tab>

          <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);">
            <Personal v-bind:UserInfo="UserInfo" v-bind:ChangedProfileInfo="ChangedProfileInfo"/>
          </v-tab-item>
          <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);">
            <Notifications v-bind:UserInfo="UserInfo"/>
          </v-tab-item>
          <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);">
            <Instance />
          </v-tab-item>
          <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);">
            <Deployment v-bind:Permissions="Permissions"/>
          </v-tab-item>
          <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);">
            <Users />
          </v-tab-item>
          <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);">
            <Ticket />
          </v-tab-item>
          <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);">
            <Integrations />
          </v-tab-item>
      </v-tabs>
  </div> 
</template>
<script lang="ts">
import Vue from "vue";
import Drawer from "../Main/Drawer.vue";
import AppBar from "../Main/AppBar.vue";
import Personal from "./Personal.vue";
import Ticket from "./Ticket.vue";
import Integrations from "./Integrations.vue";
import Notifications from "./Notifications.vue";
import Users from "./Users.vue";
import Instance from "./Instance.vue";
import Deployment from "./Deployment.vue";


interface Permissions {
  Admin: boolean;
  CanCreateUsers: boolean;
  CanModifyUsers: boolean;
  CanDeleteUsers: boolean;
  CanCreateGroups: boolean;
  CanModifyGroups: boolean;
  CanDeleteGroups: boolean;
  CanChangePermissionsGlobal: boolean;
}
let Permissions: Permissions


interface User {
ID:          number;
Username:    string;
Mail:        string;
Firstname:   string;
Lastname:    string;
Password:    string;
Settings:    {};
}

const UserInfo: User = {
    ID: 0,
    Username: "",
    Firstname: "",
    Lastname: "",
    Mail: "",
    Password: "",
    Settings: {}
}

//Initialize seperately so we don't create a reference
const ChangedProfileInfo: User = {
    ID: 0,
    Username: "",
    Firstname: "",
    Lastname: "",
    Mail: "",
    Password: "",
    Settings: {}
}

export default Vue.extend({
  name: "Settings",
  components: {
    Drawer,
    Personal,
    Ticket,
    Integrations,
    Notifications,
    Users,
    Instance,
    Deployment
  },
  data: function() {
    return {
      Permissions,
      SelectedTab: 0,
      UserInfo,
      ChangedProfileInfo
    };
  },
  mounted: async function() {
    const user = await Vue.prototype.$Request("GET", "/api/v1/profile");
    this.Permissions = await Vue.prototype.$Request("GET", "/api/v1/user/" + user.ID + "/permissions")
    this.SelectRightTab()

    const data =  await Vue.prototype.$Request("GET", "/api/v1/profile")
    //Assign so we dont create a reference here
    Object.assign(this.UserInfo, data);
     Object.assign(this.ChangedProfileInfo, data);
  },
  watch:{
      $route (to, from){
          this.SelectRightTab()
      }
  },
  methods:{
    SelectRightTab: function(){
      if(this.$route.query.setting != undefined){
        const tab = Number.parseInt(this.$route.query.setting as string)
        if (!isNaN(tab)){
          this.SelectedTab = tab
        }
      }else{
        this.UpdateQuery()
      }
    },
    UpdateQuery: function(){
      try{
        this.$router.push({ query: Object.assign({}, this.$route.query, { setting: this.SelectedTab }) });
      }finally{
        //Do nothing
      }
    }
  }
});
</script>