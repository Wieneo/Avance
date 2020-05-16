<template>
    <div >
         <v-skeleton-loader
          :loading="Loading"
          :transition="scale-transition"
          type="article"
        >
            <v-card>
                <v-card-title>{{CurrentTicket.Title}}</v-card-title>
                <v-card-subtitle>{{CurrentTicket.Description}}</v-card-subtitle>
                <v-card-text>Action!</v-card-text>
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
          //LoadTicketInformation
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
            this.CurrentTicket = (await Vue.prototype.$GetRequest("/api/v1/ticket/" + TicketID)).Ticket
            this.Loading = false
        }
    }
  })
</script>