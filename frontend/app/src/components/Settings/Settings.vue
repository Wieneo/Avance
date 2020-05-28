<template>
  <div>
    <Drawer v-on:ShowProjects="showProjects = true" />
    <div>
      <v-tabs background-color="#46494c" dark height="40px" style="margin-left: 56px;">
        <v-tab>
          <v-icon left>mdi-account</v-icon>Personal
        </v-tab>
        <v-tab>
          <v-icon left>mdi-bell-ring</v-icon>Notifications
        </v-tab>
        <v-tab v-if="UserInfo.Permissions.Admin">
            <v-icon left>mdi-cog</v-icon>Instance
        </v-tab>

        <v-tab v-if="UserInfo.Permissions.Admin">
            <v-icon left>mdi-human</v-icon>User & Groups
        </v-tab>

        <v-tab v-if="UserInfo.Permissions.Admin">
            <v-icon left>mdi-ticket</v-icon>Ticket
        </v-tab>

        <v-tab v-if="UserInfo.Permissions.Admin">
            <v-icon left>mdi-transit-connection-variant</v-icon>Integrations
        </v-tab>

        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);"><Personal/></v-tab-item>
        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);"><Notifications/></v-tab-item>
        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);"><Instance/></v-tab-item>
        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);"><Users/></v-tab-item>
        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);"><Ticket/></v-tab-item>
        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);"><Integrations/></v-tab-item>
      </v-tabs>
    </div>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Drawer from "../Main/Drawer.vue";
import AppBar from "../Main/AppBar.vue";
import Personal from "./Personal.vue";
import Ticket from "./Ticket.vue"
import Integrations from "./Integrations.vue"
import Notifications from "./Notifications.vue"
import Users from "./Users.vue"
import Instance from "./Instance.vue"

interface User {
  ID: number;
  Username: string;
  Mail: string;
  Firstname: string;
  Lastname: string;
  Password: string;
  Permissions: {};
}
const UserInfo: User = {
  ID: 0,
  Username: "Loading",
  Firstname: "Please",
  Lastname: "Wait",
  Mail: "Loading",
  Password: "",
  Permissions: {}
};

export default Vue.extend({
  name: "Settings",
  components: {
    Drawer,
    Personal,
    Ticket,
    Integrations,
    Notifications,
    Users,
    Instance
  },
  data: function() {
    return {
      UserInfo,
    };
  },
  mounted: async function() {
    const data = await Vue.prototype.$GetRequest("/api/v1/profile");
    //Assign so we dont create a reference here
    Object.assign(this.UserInfo, data);
  }
});
</script>