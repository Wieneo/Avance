<template>
    <div>
            <v-card style="height: calc(100vh - 88px); overflow-y: auto; overflow-x: hidden; max-width: 100%;" flat tile>
              <v-row style="margin: 0; padding: 0;">
                <v-col style="margin: 0; padding: 0;">
                  <div style="height: 40px; background-color: #1976d2;">&nbsp;</div>
                </v-col>
              </v-row>
              <v-skeleton-loader
                :loading="Loading"
                transition="fade-transition"
                type="article"
              >
              <v-row style="margin-top: 0px; padding-top: 0px;">
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
                  <p title="Created" class="TicketDisplayProperty"><v-icon>mdi-plus</v-icon><span class="TicketDisplayPropertyText">{{CurrentTicket.CreatedAt}}</span></p>
                  <p title="Last Modified" class="TicketDisplayProperty"><v-icon>mdi-update</v-icon><span class="TicketDisplayPropertyText">{{CurrentTicket.LastModified}}</span></p>
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
                    Recipients || PLACEHOLDER
                  </div>
                  </v-card-text>
                </v-col>
              </v-row>
              </v-skeleton-loader>
            </v-card>
    </div>
</template>
<script lang="ts">
  import Vue from 'vue'

  export default Vue.extend({
    name: 'TicketDisplay',
    mounted: async function(){
      if(this.$route.query.ticket != undefined){
        const ticketID = parseInt(this.$route.query.ticket as string)
        if (!isNaN(ticketID)){
          this.GetTicket(ticketID)
        }
      }
    },
    props: ["CurrentTicketID"],
    data: function(){
      return {
        CurrentTicket: {},
        Loading: true
      }
    },
    watch: {
        CurrentTicketID: function (val) {
            this.GetTicket(val)
        }
    },
    methods:{
        GetTicket: async function(TicketID: number){
            this.Loading = true
            this.CurrentTicket = (await Vue.prototype.$GetRequest("/api/v1/ticket/" + TicketID))
            this.Loading = false
        },
        GoToTicket: async function(TicketID: number){
          try{
            this.$router.push({ query: Object.assign({}, this.$route.query, { ticket: TicketID }) });
            this.GetTicket(TicketID)
          }finally{
            //Do nothing
          }
        }
    }
  })
</script>