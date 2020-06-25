<template>
    <div style="margin: 5px;">
        <v-skeleton-loader
          class="mx-auto"
          type="text@1, paragraph@3, text@2, paragraph@1"
          v-if="TicketLoading"
        ></v-skeleton-loader>
          <v-card v-for="action in CurrentTicket.Actions" :key="action.ID" style="margin-top: 5px;" v-else>
            <div v-if="action.Type == 0 || action.Type == 1">
              <v-row style="margin:0;">
                <v-col style="margin:0;">
                  <v-card-subtitle><b>{{action.Title}}</b><br>
                    <span v-if="action.IssuedBy.Valid">{{action.IssuedBy.Issuer.Firstname}} {{action.IssuedBy.Issuer.Lastname}} ({{action.IssuedBy.Issuer.Username}})</span>
                    <span v-else><b>System</b></span>
                    <br>{{action.IssuedAt | moment("dddd, MM/DD/YYYY HH:mm:ss")}}
                  </v-card-subtitle>
                  <v-card-text style="color: black;" v-html="action.Content"></v-card-text>
                </v-col>
                <v-col lg="1" sm="2" style="text-align: center;">
                  <v-progress-linear
                      indeterminate
                      color="primary"
                      v-if="TasksLoading"
                      style="margin-top: 10px;"
                    ></v-progress-linear>
                  <!--If no task is running for that action-->
                  <span v-else style="animation: 1s ease-out 0s 1 scaleOut;">
                    <v-btn icon v-if="!RunningTasks.get(action.ID)">
                      <v-icon v-if="!hasFailedTasks(action.ID)" style="color: #0174ff;" >mdi-check</v-icon>
                      <v-icon v-else style="color: rgba(255, 52, 52, 0.76);" >mdi-exclamation-thick</v-icon>
                    </v-btn>
                    <v-progress-circular v-else
                      style="height: 24px; cursor: pointer;"
                      indeterminate
                      color="primary"
                      @click="testLoaderClick"
                    ></v-progress-circular>
                  </span>
                </v-col>
              </v-row>
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
      resetTasks: function(){
        this.TasksLoading = true
      },
      tasksLoaded: function(){
        this.$forceUpdate();
        this.TasksLoading = false
      },
      hasFailedTasks: function(ActionID: bigint): boolean{
        let foundFaulty = false
        this.CurrentTicket.Actions.forEach(element => {
          if (element.ID == ActionID){
            element.ResolvedTasks.forEach(task => {
              if (task.Status == 2){
                foundFaulty = true
              }
            });
          }
        });

        return foundFaulty
      },
      testLoaderClick: function(){
        console.log("OK")
      }
    }
  })
</script>