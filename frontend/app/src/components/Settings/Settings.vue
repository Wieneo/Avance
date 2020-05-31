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

        <v-tab v-if="Permissions.Admin">
          <v-icon left>mdi-cog</v-icon>Instance
        </v-tab>

        <v-tab v-if="Permissions.Admin">
          <v-icon left>mdi-cog</v-icon>Deployment
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
          <Personal />
        </v-tab-item>
        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);">
          <Notifications />
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
      Permissions
    };
  },
  mounted: async function() {
    const user = await Vue.prototype.$Request("GET", "/api/v1/profile");
    this.Permissions = await Vue.prototype.$Request("GET", "/api/v1/user/" + user.ID + "/permissions")
  }
});
</script>