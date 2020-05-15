<template>
    <v-card>
        <div v-for="queue in Queues" :key="queue.ID">
            <v-list subheader dense>
                <v-subheader>{{queue.Name}}</v-subheader>

                <v-list-item
                    v-for="index in 5"
                    :key="index"
                    @click="console.log('test')">
                    <v-list-item-avatar>
                        <v-img src="https://randomuser.me/api/portraits/women/85.jpg"></v-img>
                    </v-list-item-avatar>

                    <v-list-item-content>
                        <v-list-item-title>Ticket {{index}}</v-list-item-title>
                        <v-list-item-subtitle>Open</v-list-item-subtitle>
                    </v-list-item-content>

                    <v-list-item-icon>
                        <v-icon>mdi-fire</v-icon>
                        <v-icon @click="console.log('test2')">mdi-forward</v-icon>
                    </v-list-item-icon>
                </v-list-item>
            </v-list>
            <v-divider/>
        </div>
    </v-card>
</template>
<script lang="ts">
  import Vue from 'vue'

  interface Queue{
    ID: number;
    Name: string;
  }

  const Queues: Queue[] = []

  export default Vue.extend({
    name: 'TicketList',
     data: function(){
      return {
          Queues
      }
    },
    mounted: async function(){
      this.Queues = (await Vue.prototype.$GetRequest("/api/v1/project/1/queues")).Queues
    }
  })
</script>