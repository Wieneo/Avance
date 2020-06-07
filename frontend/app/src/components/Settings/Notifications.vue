<template>
    <v-container>
        <v-row>
            <v-col>
                <v-card>
                    <v-card-title>
                        Channels
                    </v-card-title>
                    <v-card-text>
                        <v-overlay absolute=true :value="Loading">
                            <v-progress-circular indeterminate size="64"></v-progress-circular>
                        </v-overlay>
                        <v-checkbox v-model="UserInfo.Settings.EnabledNotificationChannelsReadable" label="E-Mail" value="Mail" @click="UpdateChannels"></v-checkbox>
                        <v-checkbox v-model="UserInfo.Settings.EnabledNotificationChannelsReadable" label="Telegram (Coming soon)" value="Telegram" disabled></v-checkbox>
                    </v-card-text>
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>
<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
    name: 'Notifications',
    props:["UserInfo"],
    data: function() {
        return {
            Loading: false
        }
    },
    methods:{
        UpdateChannels: async function(){
            this.Loading = true
            
            this.Loading = false
        }
    }
})
</script>