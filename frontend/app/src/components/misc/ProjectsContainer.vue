<template>
    <v-dialog
        v-model="showProjects"
        width="500"
        @input="$emit('closeProjects')"
      >
      <v-card>
          <v-card-title
            class="headline grey lighten-2"
            primary-title
          >
            Projects
          </v-card-title>

          <v-card-text>
            <v-list two-line>
              <template>
                <v-list-item v-for="project in this.Projects" :key="project.ID" @click="ChooseProject(project.ID)">
                  <v-list-item-avatar>
                    <v-img src="https://randomuser.me/api/portraits/women/85.jpg"></v-img>
                  </v-list-item-avatar>

                  <v-list-item-content>
                    <v-list-item-title>{{project.Name}}</v-list-item-title>
                    <v-list-item-subtitle>{{project.Description}}</v-list-item-subtitle>
                  </v-list-item-content>
                </v-list-item>
              </template>
            </v-list>
          </v-card-text>

          <v-divider></v-divider>

          <!--<v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              color="primary"
              text
              @click="$emit('closeProjects')"
            >
              I accept
            </v-btn>
          </v-card-actions-->
        </v-card>
      </v-dialog>
</template>

<script lang="ts">
  import Vue from 'vue'

  interface Project{
    ID: number;
    Name: string;
    Description: string;
  }

  const Projects: Project[] = []

  export default Vue.extend({
    name: 'ProjectsContainer',
    props: ['showProjects'],
    data: function(){
      return {
          Projects
      }
    },
    mounted: async function(){
      this.Projects = (await Vue.prototype.$GetRequest("/api/v1/projects"))
      //ToDo: Get to last project

      if(this.$route.name == "Main" && this.$route.query.project == undefined && this.Projects.length > 0){
        this.$router.push({ path: '/', query: { project: this.Projects[0].ID.toString() } })
      }
    },
    methods: {
      ChooseProject: function(ProjectID: number){
        try{
          this.$router.push({ path: '/', query: { project: ProjectID.toString() } })
        }finally{
          this.$emit('closeProjects')
        }
      }
    }
  })
</script>