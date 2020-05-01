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
                <v-list-item v-for="index in 20" :key="index" @click="ChooseProject(index)">
                  <v-list-item-avatar>
                    <v-img src="https://randomuser.me/api/portraits/women/85.jpg"></v-img>
                  </v-list-item-avatar>

                  <v-list-item-content>
                    <v-list-item-title>Project {{index}}</v-list-item-title>
                    <v-list-item-subtitle>This is the first project</v-list-item-subtitle>
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

  export default Vue.extend({
    name: 'ProjectsContainer',
    props: ['showProjects'],
    methods: {
      ChooseProject: function(ProjectID: number){
        //Switch Project
      }
    }
  })
</script>