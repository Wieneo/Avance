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
                    <v-btn icon v-if="!RunningTasks.get(action.ID)" @click="showTasksDetailed(action)">
                      <v-icon v-if="!hasFailedTasks(action.ID)" style="color: #0174ff;" >mdi-check</v-icon>
                      <v-icon v-else style="color: rgba(255, 52, 52, 0.76);" >mdi-exclamation-thick</v-icon>
                    </v-btn>
                    <v-progress-circular v-else
                      style="height: 24px; cursor: pointer;"
                      indeterminate
                      color="primary"
                      @click="showTasksDetailed(action)"
                    ></v-progress-circular>
                  </span>
                </v-col>
              </v-row>
            </div>
          </v-card>
          <v-dialog
            v-model="ShowTaskDetails"
            width="800"
            @input="ShowTaskDetails = false"
          >
          <v-card>
              <v-card-title
                class="headline grey lighten-2"
                primary-title
              >
                Tasks
              </v-card-title>

              <v-card-text>
                <v-list two-line>
                  <template v-if="CurrentTaskAction.Tasks != undefined">
                    <div v-if="CurrentTaskAction.Tasks.length > 0">
                    <v-list-item v-for="task in CurrentTaskAction.ResolvedTasks" :key="task.ID">
                      <v-list-item-avatar>
                        <v-progress-circular v-if="task.Status == 1 || task.Status == 0"
                          style="height: 24px; cursor: pointer;"
                          indeterminate
                          color="primary"
                          title="Task is still running"
                        ></v-progress-circular>
                         <v-icon v-if="task.Status == 3" style="color: #0174ff;" title="Task succeded">mdi-check</v-icon>
                         <v-icon v-if="task.Status == 2" style="color: rgba(255, 52, 52, 0.76);" title="Task has failed">mdi-exclamation-thick</v-icon>
                      </v-list-item-avatar>
                      <v-list-item-content>
                        <v-list-item-title>Notify {{task.Recipient.String}} about updates</v-list-item-title>
                        <v-list-item-subtitle v-if="task.Data.NotifyType == 0">Answer was added</v-list-item-subtitle>
                        <v-list-item-subtitle v-if="task.Data.NotifyType == 1">Comment was added</v-list-item-subtitle>

                        <v-divider></v-divider>
                        <v-expansion-panels accordion flat hover>
                            <v-expansion-panel>
                              <v-expansion-panel-header>Results</v-expansion-panel-header>
                              <v-expansion-panel-content>
                                <ul>
                                  <li>[{{task.QueuedAt | moment("dddd, MM/DD/YYYY HH:mm:ss")}}] Job queued</li>
                                  <li v-for="result in task.Results" :key="result.Result">[{{result.IssuedAt | moment("dddd, MM/DD/YYYY HH:mm:ss")}}] {{result.Result}}</li>
                                </ul>
                              </v-expansion-panel-content>
                            </v-expansion-panel>
                          </v-expansion-panels>
                          <v-divider></v-divider>
                      </v-list-item-content>

                    </v-list-item>
                    </div>
                    <div v-else>
                      <v-list-item-subtitle>There are no tasks assigned to that action.</v-list-item-subtitle>
                    </div>
                  </template>
                </v-list>
              </v-card-text>

              <v-divider></v-divider>
            </v-card>
          </v-dialog>
    </div>
</template>
<script lang="ts">
  import Vue from 'vue'

   export default Vue.extend({
    name: 'ActionDisplay',
    props: ["CurrentTicket", "TicketLoading"],
    data: function(){
        return {
            TasksLoading: true,
            CurrentTaskAction: {},
            ShowTaskDetails: false
        }
    },
    computed:{
      RunningTasks:{
        cache: false,
        get (){
          const actionTasks = new Map<bigint, boolean>()
          this.CurrentTicket.Actions.forEach((element: any) => {
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
        this.CurrentTicket.Actions.forEach((element: any) => {
          if (element.ID == ActionID){
            element.ResolvedTasks.forEach((task: any) => {
              if (task.Status == 2){
                foundFaulty = true
              }
            });
          }
        });

        return foundFaulty
      },
      showTasksDetailed: function(Action: any){
        this.CurrentTaskAction = Action;
        this.ShowTaskDetails = true
      }
    }
  })
</script>