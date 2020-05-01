<template>
<v-app>
  <Login v-if="this.$route.name === 'Login'"/>
  <MainApp v-if="this.$route.name === 'Main'"/>
</v-app>
</template>

<script lang="ts">
import Vue from 'vue';
import MainApp from './components/Main/MainApp.vue'
import Login from './components/Login/Login.vue'


export default Vue.extend({
  name: 'App',
  components: {
    MainApp,
    Login
  },
});
</script>
