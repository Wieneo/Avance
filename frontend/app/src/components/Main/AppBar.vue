<template>
    <v-app-bar app>
        <p>Project</p>
        <v-spacer></v-spacer>
        <v-btn icon>
            <v-icon>mdi-logout</v-icon>
        </v-btn>
    </v-app-bar>
</template>

<script lang="ts">
  import Vue from 'vue'

  export default Vue.extend({
    name: 'AppBar',
  })
</script>