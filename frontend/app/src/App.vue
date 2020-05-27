<template>
<v-app>
  <Login v-if="this.$route.name === 'Login'"/>
  <MainApp v-if="this.$route.name === 'Main'"/>
  <Settings v-if="this.$route.name === 'Settings'"/>
</v-app>
</template>

<style scoped>
@import '../node_modules/noty/lib/noty.css';
@import '../node_modules/noty/lib/themes/metroui.css';
@import '../public/custom.css';
</style>

<script lang="ts">
import Vue from 'vue';
import MainApp from './components/Main/MainApp.vue'
import Login from './components/Login/Login.vue'
import Settings from './components/Settings/Settings.vue'

export default Vue.extend({
  name: 'App',
  components: {
    MainApp,
    Login,
    Settings
  }
});
</script>
