<template>
    <div>
      <v-skeleton-loader
          v-if="LoadingTickets"
          type="list-item-avatar-two-line@5"
        ></v-skeleton-loader>
        <v-card v-else>
            <div v-for="queue in Queues" :key="queue.ID">
                <v-list subheader dense>
                    <v-subheader>{{queue.Name}}</v-subheader>
                    <v-list-item-group v-model="queue.SelectedTicket">
                      <v-list-item
                          v-for="ticket in queue.Tickets"
                          :key="ticket.ID"
                          @click="DisplayTicket(ticket.ID)">
                          <v-list-item-avatar>
                              <v-img v-if="ticket.OwnerID.Valid" :src="getUserAvatarLink(ticket.Owner.ID)"></v-img>
                              <v-img v-else src="" style="background-color: #d0d0d0;"></v-img>
                          </v-list-item-avatar>

                          <v-list-item-content>
                              <v-list-item-title>{{ticket.Title}}</v-list-item-title>
                              <v-list-item-subtitle :style="{ color: ticket.Status.DisplayColor }">{{ticket.Status.Name}}</v-list-item-subtitle>
                          </v-list-item-content>

                          <v-list-item-icon>
                              <v-icon :style="{ color: ticket.Severity.DisplayColor }" :title="ticket.Severity.Name">mdi-fire</v-icon>
                          </v-list-item-icon>
                      </v-list-item>
                    </v-list-item-group>
                </v-list>
                <v-divider/>
            </div>
        </v-card>
    </div>
</template>
<script lang="ts">
  import Vue from 'vue'

  interface Queue{
    ID: number;
    Name: string;
    Tickets: any[];
  }

  const Queues: Queue[] = []

  export default Vue.extend({
    name: 'TicketList',
    watch:{
        $route (to, from){
            if(to.query.project != undefined){
              const projectID = parseInt(to.query.project  as string)
              if (!isNaN(projectID)){
                if (this.CurrentProject != projectID){
                  this.CurrentProject = projectID
                  this.LoadQueues()
                }
              }
            }
        }
    },
     data: function(){
      return {
        LoadingTickets: true,
        CurrentProject: 0,
        Queues
      }
    },
    mounted: async function(){
      if(this.$route.query.project != undefined){
        const projectID = parseInt(this.$route.query.project as string)
        if (!isNaN(projectID)){
          if (this.CurrentProject != projectID){
            this.CurrentProject = projectID
            this.LoadQueues()
          }
        }
      }
    },
    methods:{
      LoadQueues: async function(){
        this.LoadingTickets = true
        const rawQueues = (await Vue.prototype.$Request("GET", "/api/v1/project/" + this.CurrentProject + "/queues"))

        await this.asyncForEach(rawQueues, async (element: any, index: number) => {
          rawQueues[index].Tickets = (await Vue.prototype.$Request("GET", "/api/v1/project/" + this.CurrentProject + "/queue/" + element.ID + "/tickets"))
        });

        //Fix for v-list selection
         if(this.$route.query.ticket != undefined){
          const ticketID = parseInt(this.$route.query.ticket as string)
          if (!isNaN(ticketID)){
            rawQueues.forEach((queue: any, index: number) => {
              queue.Tickets.forEach((element: any, tindex: number) => {
                if (element.ID == ticketID){
                  rawQueues[index].SelectedTicket = tindex
                }
              });
            });
          }
        }

        this.Queues = rawQueues

        this.LoadingTickets = false
      },
      DisplayTicket: async function(TicketID: number){
        //Fix for v-list selection
        this.Queues.forEach((queue: any, index: number) => {
          queue.Tickets.forEach((element: any, tindex: number) => {
            if (element.ID == TicketID){
              this.Queues[index].SelectedTicket = tindex
            }else{
              this.Queues[index].SelectedTicket = -1
            }
          });
        });

        try{
          this.$router.push({ query: Object.assign({}, this.$route.query, { ticket: TicketID }) });
        }finally{
          //Do Nothing
        }

        this.$forceUpdate();
      },
      asyncForEach: async function (array: any, callback: any) {
        for (let index = 0; index < array.length; index++) {
          await callback(array[index], index, array);
        }
      },
      getUserAvatarLink(UserID: number){
        return '/api/v1/user/' + UserID + '/avatar?' + performance.now()
      }
    }
  })
</script>