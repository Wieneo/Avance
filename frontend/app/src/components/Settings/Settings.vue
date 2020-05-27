<template>
  <div>
    <Drawer v-on:ShowProjects="showProjects = true" />
    <div>
      <v-tabs background-color="primary" dark height="40px" style="margin-left: 56px;">
        <v-tab>
          <v-icon left>mdi-account</v-icon>Personal
        </v-tab>
        <v-tab>
          <v-icon left>mdi-bell-ring</v-icon>Notifications
        </v-tab>
        <v-tab v-if="UserInfo.Permissions.Admin">
            <v-icon left>mdi-cog</v-icon>Instance
        </v-tab>

        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);">Item1</v-tab-item>
        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);">Item2</v-tab-item>
      </v-tabs>
    </div>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Drawer from "../Main/Drawer.vue";
import AppBar from "../Main/AppBar.vue";

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
    Drawer
  },
  data: function() {
    return {
      UserInfo
    };
  },
  mounted: async function() {
    const data = await Vue.prototype.$GetRequest("/api/v1/profile");
    //Assign so we dont create a reference here
    Object.assign(this.UserInfo, data);
  }
});
</script>