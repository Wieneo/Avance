<template>
    <div>
        <Drawer v-on:ShowProjects="showProjects = true"/>
        <AppBar/>
        <v-content>

            <!-- Provides the application the proper gutter -->
            <v-container fluid style="max-height: calc(100vh - 64px); overflow-y: auto">
            <ProjectsContainer v-bind:showProjects="showProjects" v-on:closeProjects="showProjects = false"/>
            <v-row no-gutters>
                <v-col lg="3">
                <TicketList style="max-height: calc(100vh - 88px); overflow-y: auto" v-on:showTicket="DisplayTicket"/>
                </v-col>
                <v-col>
                     <v-tabs
                        v-model="tab"
                        background-color="primary"
                        dark
                        height="40px"
                    >
                        <v-tab><v-icon left>mdi-account</v-icon>Interactions</v-tab>
                        <v-tab><v-icon left>mdi-history</v-icon>Actions</v-tab>

                        <v-tab-item>Normal Actions</v-tab-item>
                        <v-tab-item>All Actions</v-tab-item>
                    </v-tabs>
                </v-col>
                <v-col lg="2" v-if="CurrentTicketID != 0">
                   <TicketDisplay v-bind:CurrentTicketID="CurrentTicketID"/>
                </v-col>
                <v-col v-else>
                    <!--ToDo: Placeholder if no ticket is selected-->
                </v-col>
            </v-row>
            </v-container>
        </v-content>
    </div>
</template>

<script lang="ts">
import Vue from 'vue'
import Drawer from './Drawer.vue';
import AppBar from './AppBar.vue';
import ProjectsContainer from '../misc/ProjectsContainer.vue';
import TicketList from './TicketList.vue';
import TicketDisplay from './TicketDisplay.vue';
//import ActionDisplay from './ActionDisplay.vue';


export default Vue.extend({
    name: "MainApp",
    components:{
        Drawer,
        AppBar,
        ProjectsContainer,
        TicketList,
        TicketDisplay,
    },
    data: function(){
        return {
            showProjects: false,
            CurrentTicketID: 0,
        }
    },
    mounted: async function(){
      if(this.$route.query.ticket != undefined){
        const ticketID = parseInt(this.$route.query.ticket as string)
        if (!isNaN(ticketID)){
          this.DisplayTicket(ticketID)
        }
      }
    },
    methods:{
        DisplayTicket: function(TicketID: number){
            this.CurrentTicketID = TicketID
        }
    }
})
</script>