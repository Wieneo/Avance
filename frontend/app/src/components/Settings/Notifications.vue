<template>
    <v-container>
        <v-row>
            <v-col>
                <v-card>
                    <v-card-title>
                        Channels
                    </v-card-title>
                    <v-card-text>
                        <v-overlay :absolute=true :value="ChannelsLoading">
                            <v-progress-circular indeterminate size="64"></v-progress-circular>
                        </v-overlay>
                        <v-checkbox v-model="UserInfo.Settings.Notification.MailNotificationEnabled" label="E-Mail" @change="UpdateChannels"></v-checkbox>
                        <v-checkbox v-model="UserInfo.Settings.Notification.TelegramNotificationEnabled" label="Telegram (Coming soon)" value="Telegram" disabled></v-checkbox>
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
            ChannelsLoading: false
        }
    },
    methods:{
        UpdateChannels: async function(){
            this.ChannelsLoading = true
            if (!(await this.Update())){
                this.$emit('refreshUserInfo')
            }
            this.ChannelsLoading = false
        },
        Update: async function(): Promise<boolean>{
            const result = await Vue.prototype.$Request("PATCH", "/api/v1/profile/settings", this.UserInfo.Settings)
            if (result.Error == undefined){
                return true
            }
            return false
        }
    }
})
</script>