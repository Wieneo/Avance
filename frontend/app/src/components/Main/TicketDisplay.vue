<template>
    <div>
         <v-skeleton-loader
          :loading="Loading"
          transition="fade-transition"
          type="article"
        >
            <v-card>
              <v-row>
                <v-col>
                  <v-card-title>{{CurrentTicket.Title}}</v-card-title>
                  <v-card-subtitle>{{CurrentTicket.Description}}</v-card-subtitle>
                </v-col>
                <v-col lg="3">
                  <p title="ID" class="TicketDisplayProperty"><v-icon>mdi-pound-box-outline</v-icon><span class="TicketDisplayPropertyText">{{CurrentTicket.ID}}</span></p>
                  <p title="Queue" class="TicketDisplayProperty"><v-icon>mdi-tray-full</v-icon><span class="TicketDisplayPropertyText">{{CurrentTicket.Queue.Name}}</span></p>
                  <p title="Status" class="TicketDisplayProperty"><v-icon>mdi-circle-outline</v-icon><span class="TicketDisplayPropertyText" :style="{ color: CurrentTicket.Status.DisplayColor }">{{CurrentTicket.Status.Name}}</span></p>
                  <p title="Severity" class="TicketDisplayProperty"><v-icon>mdi-fire</v-icon><span class="TicketDisplayPropertyText" :style="{ color: CurrentTicket.Severity.DisplayColor }">{{CurrentTicket.Severity.Name}}</span></p>
                </v-col>
                <v-col lg="3">
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
                </v-col>
              </v-row>
            </v-card>
         </v-skeleton-loader>
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
        }
    }
  })
</script>