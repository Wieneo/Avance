<template>
    <v-app-bar app>
        <p>Project</p>
        <v-spacer></v-spacer>
        <v-btn icon>
            <v-icon @click="logout">mdi-logout</v-icon>
        </v-btn>
    </v-app-bar>
</template>

<script lang="ts">
  import Vue from 'vue'

  export default Vue.extend({
    name: 'AppBar',
    methods:{
      logout: async function(){
        await Vue.prototype.$GetRequest("/api/v1/logout")
        window.location.href = "/login"
      }
    }
  })
</script>