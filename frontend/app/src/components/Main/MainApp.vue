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
                <TicketList style="max-height: calc(100vh - 88px); overflow-y: auto"/>
                </v-col>
                <v-col>
                Ticket
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

export default Vue.extend({
    name: "MainApp",
    components:{
        Drawer,
        AppBar,
        ProjectsContainer,
        TicketList
    },
    data: function(){
        return {
        showProjects: false
        }
    }
})
</script>