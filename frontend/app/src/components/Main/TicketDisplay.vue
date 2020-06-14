<template>
    <div>
            <v-card style="height: calc(100vh - 88px); overflow-y: auto; overflow-x: hidden; max-width: 100%;" flat tile>
              <v-row style="margin: 0; padding: 0;">
                <v-col style="margin: 0; padding: 0;">
                  <div style="height: 40px; background-color: #1976d2;">&nbsp;</div>
                </v-col>
              </v-row>
              <v-skeleton-loader v-if="TicketLoading"
                transition="fade-transition"
                type="article"
              ></v-skeleton-loader>
              <v-row style="margin-top: 0px; padding-top: 0px; max-height: calc(100vh - 88px);" class="overflow-y-auto" v-else>
                <v-col style="margin-top: 0px; padding-top: 0px;">
                  <v-card-title>{{CurrentTicket.Title}}</v-card-title>
                  <v-card-subtitle>{{CurrentTicket.Description}}</v-card-subtitle>
                  <v-card-text>
                  <p title="ID" class="TicketDisplayProperty"><v-icon>mdi-pound-box-outline</v-icon><span class="TicketDisplayPropertyText">{{CurrentTicket.ID}}</span></p>
                  <p title="Queue" class="TicketDisplayProperty"><v-icon>mdi-tray-full</v-icon><span class="TicketDisplayPropertyText">{{CurrentTicket.Queue.Name}}</span></p>
                  <p title="Status" class="TicketDisplayProperty"><v-icon>mdi-circle-outline</v-icon><span class="TicketDisplayPropertyText" :style="{ color: CurrentTicket.Status.DisplayColor }">{{CurrentTicket.Status.Name}}</span></p>
                  <p title="Severity" class="TicketDisplayProperty"><v-icon>mdi-fire</v-icon><span class="TicketDisplayPropertyText" :style="{ color: CurrentTicket.Severity.DisplayColor }">{{CurrentTicket.Severity.Name}}</span></p>

                  <p title="Owner" class="TicketDisplayProperty">
                    <v-icon>mdi-account-circle-outline</v-icon>
                    <span v-if="CurrentTicket.OwnerID.Valid" class="TicketDisplayPropertyText">{{CurrentTicket.Owner.Username}} ({{CurrentTicket.Owner.Firstname}} {{CurrentTicket.Owner.Lastname}})</span>
                    <span v-else class="TicketDisplayPropertyText">Nobody</span>
                  </p>
                  <p title="Created" class="TicketDisplayProperty"><v-icon>mdi-plus</v-icon><span class="TicketDisplayPropertyText">{{CurrentTicket.CreatedAt | moment("dddd, MM/DD/YYYY HH:mm:ss")}}</span></p>
                  <p title="Last Modified" class="TicketDisplayProperty"><v-icon>mdi-update</v-icon><span class="TicketDisplayPropertyText">{{CurrentTicket.LastModified | moment("dddd, MM/DD/YYYY HH:mm:ss")}}</span></p>
                  <p title="Stalled Until" class="TicketDisplayProperty">
                    <v-icon>mdi-calendar-range</v-icon>
                    <span v-if="CurrentTicket.StalledUntil.Valid" class="TicketDisplayPropertyText">{{CurrentTicket.StalledUntil.Time}}</span>
                    <span v-else class="TicketDisplayPropertyText">None</span>
                  </p>
                  <div style="margin-top: 10px;" v-if="CurrentTicket.Relations.length > 0">
                    <hr style="margin-bottom: 3px;">
                    <div v-for="relation in CurrentTicket.Relations" :key="relation.ID" @click="GoToTicket(relation.OtherTicket.ID)">
                      <span class="TicketDisplayProperty" v-if="relation.Type == 0" title="References"><v-icon>mdi-arrow-right</v-icon></span>
                      <span class="TicketDisplayProperty" v-if="relation.Type == 1" title="Referenced By"><v-icon>mdi-arrow-left</v-icon></span>
                      <span class="TicketDisplayProperty" v-if="relation.Type == 2" title="Parent of"><v-icon>mdi-human-female</v-icon></span>
                      <span  class="TicketDisplayProperty" v-if="relation.Type == 3" title="Child of"><v-icon>mdi-human-child</v-icon></span>

                      <v-menu :offset-x=true :offset-y=true :open-on-hover=true :nudge-width="200">
                        <template v-slot:activator="{ on }">
                          <a v-on="on">{{relation.OtherTicket.Title}}</a>
                        </template>
                        <v-card style="text-align: center; overflow-y: hidden;">
                          <v-card-title primary-title class="justify-center">{{relation.OtherTicket.Title}}</v-card-title>
                          <v-card-subtitle>{{relation.OtherTicket.Description}}</v-card-subtitle>
                          <p title="ID" class="TicketDisplayProperty"><v-icon>mdi-pound-box-outline</v-icon><span class="TicketDisplayPropertyText">{{relation.OtherTicket.ID}}</span></p>
                          <p title="Queue" class="TicketDisplayProperty"><v-icon>mdi-tray-full</v-icon><span class="TicketDisplayPropertyText">{{relation.OtherTicket.Queue.Name}}</span></p>
                          <p title="Status" class="TicketDisplayProperty"><v-icon>mdi-circle-outline</v-icon><span class="TicketDisplayPropertyText" :style="{ color: relation.OtherTicket.Status.DisplayColor }">{{relation.OtherTicket.Status.Name}}</span></p>
                          <p title="Severity" class="TicketDisplayProperty"><v-icon>mdi-fire</v-icon><span class="TicketDisplayPropertyText" :style="{ color: relation.OtherTicket.Severity.DisplayColor }">{{relation.OtherTicket.Severity.Name}}</span></p>
                          <p title="Owner" class="TicketDisplayProperty">
                            <v-icon>mdi-account-circle-outline</v-icon>
                            <span v-if="relation.OtherTicket.OwnerID.Valid" class="TicketDisplayPropertyText">{{relation.OtherTicket.Owner.Username}} ({{relation.OtherTicket.Owner.Firstname}} {{relation.OtherTicket.Owner.Lastname}})</span>
                            <span v-else class="TicketDisplayPropertyText">Nobody</span>
                          </p>
                        </v-card>
                      </v-menu>

                    </div>
                  </div>
                  <div style="margin-top: 10px;">
                    <hr style="margin-bottom: 3px;">
                    <span v-for="req in AllRecipients" :key="req.ID">
                      <span class="TicketDisplayProperty" title="Requestor" v-if="req.Type == 0"><v-icon>mdi-human-child</v-icon></span>
                      <span class="TicketDisplayProperty" title="Reader" v-if="req.Type == 1"><v-icon>mdi-magnify</v-icon></span>
                      <span class="TicketDisplayProperty" title="Admin" v-if="req.Type == 2"><v-icon>mdi-head-minus-outline</v-icon></span>
                      <i v-if="!req.User.Valid" style="margin-left: 10px;">{{req.Mail}}</i>
                      <v-menu :offset-x=true :offset-y=true :open-on-hover=true :nudge-width="200" v-else>
                        <template v-slot:activator="{ on }">
                          <a style="margin-left: 10px;" v-on="on">{{req.User.Value.Username}}</a>
                        </template>
                        <v-card style="text-align: center; overflow-y: hidden;">
                          <v-card-title primary-title class="justify-center">
                            <v-list-item-avatar>
                              <v-img :src="getUserAvatarLink(req.User.Value.ID)"></v-img>
                            </v-list-item-avatar>
                            {{req.User.Value.Username}}
                          </v-card-title>
                          <v-card-text>
                            {{req.User.Value.Firstname}} {{req.User.Value.Lastname}}
                            <br>
                            Status: <span style="color: green;">TBI</span><br>
                            Last Seen: <span style="color: green;">TBI</span>
                          </v-card-text>
                        </v-card>
                      </v-menu>
                      <br>
                    </span>
                  </div>
                  </v-card-text>
                </v-col>
              </v-row>
            </v-card>
    </div>
</template>
<script lang="ts">
  import Vue from 'vue'

  export default Vue.extend({
    name: 'TicketDisplay',
    props: ["CurrentTicket", "TicketLoading"],
    computed:{
      AllRecipients(){
        const data: any[] = []
        this.CurrentTicket.Recipients.Requestors.forEach((element: any) => {
          element.Type = 0
          data.push(element)
        });
        this.CurrentTicket.Recipients.Readers.forEach((element: any) => {
          element.Type = 1
          data.push(element)
        });
        this.CurrentTicket.Recipients.Admins.forEach((element: any) => {
          element.Type = 2
          data.push(element)
        });
        return data
      }
    },
    methods:{
      GoToTicket: async function(TicketID: number){
          try{
            this.$router.push({ query: Object.assign({}, this.$route.query, { ticket: TicketID }) });
          }finally{
            //Do nothing
          }
        },
      getUserAvatarLink(UserID: number){
        return '/api/v1/user/' + UserID + '/avatar?' + performance.now()
      }
    }
  })
</script>