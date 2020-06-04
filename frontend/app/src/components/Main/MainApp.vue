<template>
    <div>
        <Drawer v-on:ShowProjects="showProjects = true"/>
        <AppBar/>
        <v-content>

            <!-- Provides the application the proper gutter -->
            <v-container fluid style="max-height: calc(100vh - 64px); overflow-y: auto">
            <v-row no-gutters>
                <v-col lg="3">
                <TicketList style="max-height: calc(100vh - 88px); overflow-y: auto"/>
                </v-col>
                <v-col v-if="CurrentTicketID != 0">
                     <v-tabs
                        background-color="primary"
                        dark
                        height="40px"
                    >
                        <v-tab><v-icon left>mdi-account</v-icon>Interactions</v-tab>
                        <v-tab><v-icon left>mdi-history</v-icon>Actions</v-tab>

                        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);"><ActionDisplay v-bind:CurrentTicket="CurrentTicket" v-bind:TicketLoading="TicketLoading"/></v-tab-item>
                        <v-tab-item class="overflow-y-auto" style="max-height: calc(100vh - 130px);"><TimelineDisplay v-bind:CurrentTicket="CurrentTicket" v-bind:TicketLoading="TicketLoading"/></v-tab-item>
                    </v-tabs>
                </v-col>
                <v-col lg="2" v-if="CurrentTicketID != 0">
                   <TicketDisplay v-bind:CurrentTicket="CurrentTicket" v-bind:TicketLoading="TicketLoading"/>
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
import TicketList from './TicketList.vue';
import TicketDisplay from './TicketDisplay.vue';
import ActionDisplay from './ActionDisplay.vue';
import TimelineDisplay from './TimelineDisplay.vue';


export default Vue.extend({
    name: "MainApp",
    components:{
        Drawer,
        AppBar,
        TicketList,
        TicketDisplay,
        ActionDisplay,
        TimelineDisplay
    },
    data: function(){
        return {
            CurrentTicket: {},
            CurrentTicketID: 0,
            CurrentProjectID: 0,
            TicketLoading: false
        }
    },
    mounted: async function(){
      this.HandleRouteChange()
    },
    watch:{
        $route (to, from){
            this.HandleRouteChange()
        }
    },
    methods:{
        HandleRouteChange: function(){
            let ticketChanged = false
            if(this.$route.query.ticket != undefined){
                const ticketID = parseInt(this.$route.query.ticket as string)
                if (!isNaN(ticketID)){
                    if (this.CurrentTicketID != ticketID){
                        this.GetTicket(ticketID)
                        this.CurrentTicketID = ticketID
                        ticketChanged = true
                    }
                }
            }

            if(this.$route.query.project != undefined){
                const projectID = parseInt(this.$route.query.project as string)
                if (!isNaN(projectID)){
                    if (this.CurrentProjectID != projectID){
                        this.CurrentProjectID = projectID
                                                //This is to prevent issues on initial loading
                        if (!ticketChanged){
                            this.CurrentTicketID = 0
                            this.CurrentTicket = {}
                        }
                    }
                }
            }
        },
        GetTicket: async function(TicketID: number){
            this.TicketLoading = true
            this.CurrentTicket = (await Vue.prototype.$Request("GET", "/api/v1/ticket/" + TicketID))
            this.TicketLoading = false
        }
    }
})
</script>