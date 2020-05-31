<template>
    <v-container>
        <v-row>
            <v-col>
                <h2 style="margin-bottom: 10px;">Worker</h2>
                 <v-data-table
                    v-model="WorkersSelected"
                    :headers="headers"
                    :items="Workers"
                    item-key="name"
                    show-select
                    class="elevation-1"
                >
                <template v-slot:item.LastSeen="{ item }">
                    {{ item.LastSeen | moment("MM/DD/YYYY HH:mm:ss") }}
                </template>
                <template v-slot:item.Active="{ item }">
                    <v-btn icon v-if="item.Active" title="Active"><v-icon>mdi-check-circle-outline</v-icon></v-btn>
                    <v-btn v-else title="Inactive"><v-icon>mdi-bed</v-icon></v-btn>
                </template>
                </v-data-table>
                <v-btn color="primary" class="mr-2" style="margin-top: 10px;">Toggle Active</v-btn>
            </v-col>
        </v-row>
    </v-container>
</template>
<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
    name: 'Deployment',
    props: ["Permissions"],
    mounted: async function(){
        if (this.Permissions.Admin || this.Permissions.CanSeeWorker){
            this.Workers = await Vue.prototype.$Request("GET", "/api/v1/workers")
        }else{
            this.WorkerPermsFailed = true
        }
    },
    data: function() {
        return {
            Workers: [],
            WorkerPermsFailed: false,
            WorkersSelected: [],
            headers: [
                { text: 'ID', align: 'start', value: 'ID',},
                { text: 'Name', value: 'Name' },
                { text: 'Last Seen at', value: 'LastSeen' },
                { text: 'Active', value: 'Active' },
            ],
        };
    },
})
</script>