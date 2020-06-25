<template>
    <div style="margin: 5px;">
        <v-skeleton-loader
          class="mx-auto"
          type="text@1, paragraph@3, text@2, paragraph@1"
          v-if="TicketLoading"
        ></v-skeleton-loader>
          <v-card v-for="action in CurrentTicket.Actions" :key="action.ID" style="margin-top: 5px;" v-else>
            <div v-if="action.Type == 0 || action.Type == 1">
              <v-card-subtitle><b>{{action.Title}}</b><br>
                <span v-if="action.IssuedBy.Valid">{{action.IssuedBy.Issuer.Firstname}} {{action.IssuedBy.Issuer.Lastname}} ({{action.IssuedBy.Issuer.Username}})</span>
                <span v-else><b>System</b></span>
                <br>{{action.IssuedAt | moment("dddd, MM/DD/YYYY HH:mm:ss")}}
                <p v-if="TasksLoading">Tasks loading...</p>
                <p v-if="RunningTasks.get(action.ID) && !TasksLoading">Test</p>
              </v-card-subtitle>
              <v-card-text style="color: black;" v-html="action.Content"></v-card-text>
            </div>
          </v-card>
    </div>
</template>
<script lang="ts">
  import Vue from 'vue'

   export default Vue.extend({
    name: 'ActionDisplay',
    props: ["CurrentTicket", "TicketLoading"],
    data: function(){
        return {
            TasksLoading: true
        }
    },
    computed:{
      RunningTasks:{
        cache: false,
        get (){
          const actionTasks = new Map<bigint, boolean>()
          this.CurrentTicket.Actions.forEach(element => {
            actionTasks.set(element.ID, element.TaskRunning)
          });

          return actionTasks
        }
      }
    },
    methods:{
      tasksLoaded: function(){
        this.$forceUpdate();
        this.TasksLoading = false
      }
    }
  })
</script>