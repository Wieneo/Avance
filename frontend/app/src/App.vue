<template>
<v-app>
  <Drawer v-on:ShowProjects="showProjects = true"/>
  <AppBar/>
  <v-content>

    <!-- Provides the application the proper gutter -->
    <v-container fluid>
      <ProjectsContainer v-bind:showProjects="showProjects"/>
    </v-container>
  </v-content>
</v-app>
</template>

<script lang="ts">
import Vue from 'vue';
import Drawer from './components/Drawer.vue';
import AppBar from './components/AppBar.vue';
import ProjectsContainer from './components/misc/ProjectsContainer.vue';

export default Vue.extend({
  name: 'App',

  components: {
    Drawer,
    AppBar,
    ProjectsContainer
  },
  data: function(){
    return {
      showProjects: false
    }
  }
});
</script>
