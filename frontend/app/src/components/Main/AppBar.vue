<template>
    <v-app-bar app>
        <p>{{CurrentProject.Name}}</p>
        <v-spacer></v-spacer>
        <v-btn icon @click="logout">
            <v-icon>mdi-logout</v-icon>
        </v-btn>
    </v-app-bar>
</template>

<script lang="ts">
  import Vue from 'vue'

  export default Vue.extend({
    name: 'AppBar',
    methods:{
      logout: async function(){
        await Vue.prototype.$Request("GET", "/api/v1/logout")
        window.location.href = "/login"
      },
      getProjectInfo: async function(ProjectID: number){
          if (!isNaN(ProjectID)){
            if (this.CurrentProjectID != ProjectID){
              this.CurrentProject = (await Vue.prototype.$Request("GET", "/api/v1/project/" + ProjectID))
              this.CurrentProjectID = ProjectID
            }
          }
      }
    },
    data: function(){
      return {
        CurrentProject: {},
        CurrentProjectID: 0
      }
    },
    watch:{
        $route (to, from){
            if(to.query.project != undefined){
              this.getProjectInfo(parseInt(to.query.project  as string))
            }
        }
    },
     mounted: async function(){
      if(this.$route.query.project != undefined){
        this.getProjectInfo(parseInt(this.$route.query.project  as string))
      }
    }
  })
</script>