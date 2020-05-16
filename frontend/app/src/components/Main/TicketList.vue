<template>
    <div>
      <v-skeleton-loader
          :loading="LoadingTickets"
          :transition="scale-transition"
          type="list-item-avatar-two-line@4"
        >
        <v-card>
            <div v-for="queue in Queues" :key="queue.ID">
                <v-list subheader dense>
                    <v-subheader>{{queue.Name}}</v-subheader>
                    <v-list-item
                        v-for="ticket in queue.Tickets"
                        :key="ticket.ID"
                        @click="DisplayTicket(ticket.ID)">
                        <v-list-item-avatar>
                            <v-img src="https://randomuser.me/api/portraits/women/85.jpg"></v-img>
                        </v-list-item-avatar>

                        <v-list-item-content>
                            <v-list-item-title>{{ticket.Title}}</v-list-item-title>
                            <v-list-item-subtitle :style="{ color: ticket.Status.DisplayColor }">{{ticket.Status.Name}}</v-list-item-subtitle>
                        </v-list-item-content>

                        <v-list-item-icon>
                            <v-icon :style="{ color: ticket.Severity.DisplayColor }" :title="ticket.Severity.Name">mdi-fire</v-icon>
                            <v-icon @click="console.log('test2')">mdi-forward</v-icon>
                        </v-list-item-icon>
                    </v-list-item>
                </v-list>
                <v-divider/>
            </div>
        </v-card>
      </v-skeleton-loader>
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
        this.Queues = (await Vue.prototype.$GetRequest("/api/v1/project/" + this.CurrentProject + "/queues")).Queues
        this.LoadTickets()
      },
      LoadTickets: async function(){
        await this.asyncForEach(this.Queues, async (element: any) => {
          element.Tickets = (await Vue.prototype.$GetRequest("/api/v1/project/" + this.CurrentProject + "/queue/" + element.ID + "/tickets")).Tickets
        });

        this.LoadingTickets = false
      },
      DisplayTicket: async function(TicketID: number){
        try{
          this.$router.push({ query: Object.assign({}, this.$route.query, { ticket: TicketID }) });
        }finally{
          this.$emit('showTicket', TicketID)
        }
      },
      asyncForEach: async function (array: any, callback: any) {
        for (let index = 0; index < array.length; index++) {
          await callback(array[index], index, array);
        }
      }
    }
  })
</script>