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
                <v-card style="margin-top: 20px">
                    <v-card-title>
                        Mail Settings
                    </v-card-title>
                    <v-card-text>
                        <v-overlay :absolute=true :value="!UserInfo.Settings.Notification.MailNotificationEnabled">
                        </v-overlay>
                        <v-overlay :absolute=true :value="MailSettingsLoading">
                            <v-progress-circular indeterminate size="64"></v-progress-circular>
                        </v-overlay>
                        <v-text-field
                            v-model="UserInfo.Settings.Notification.NotificationFrequency"
                            :rules="frequencyRules"
                            label="Frequency (Seconds)"
                            required
                            @keyup="StartMailSettingsTimer"
                            @keydown="ResetMailSettingsTimer"
                            type="number"
                        ></v-text-field>
                        <v-checkbox v-model="UserInfo.Settings.Notification.NotificationAboutNewTickets" label="Notify about new tickets" @change="UpdateMailSettings"></v-checkbox>
                        <v-checkbox v-model="UserInfo.Settings.Notification.NotificationAboutUpdates" label="Notify about updates in tickets" @change="UpdateMailSettings"></v-checkbox>
                        <v-checkbox v-model="UserInfo.Settings.Notification.NotificationAfterInvolvment" label="Notify about updates in tickets you have been involved in" @change="UpdateMailSettings"></v-checkbox>
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
    props: ["UserInfo"],
    data: function() {
        return {
            ChannelsLoading: false,
            MailSettingsLoading: false,
            frequencyRules: [
                v => !!v || 'Frequency is required',
                v => /^[0-9]*$/.test(v) || 'Frequency must be a positive number',
            ],
            MailSettingsTypingTimeout: 0
        }
    },
    methods:{
        UpdateChannels: async function(){
            this.ChannelsLoading = true
            if (!(await this.Update())){
                this.$emit('refreshUserInfo')
            }else{
                Vue.prototype.$NotifySuccess("Settings updated")
            }
            this.ChannelsLoading = false
        },
        ResetMailSettingsTimer: function(){
            clearTimeout(this.MailSettingsTypingTimeout);
        },
        StartMailSettingsTimer: function(){
            clearTimeout(this.MailSettingsTypingTimeout);
            this.MailSettingsTypingTimeout = setTimeout(this.UpdateMailSettings, 2000);
        },
        UpdateMailSettings: async function(){
            this.MailSettingsLoading = true
            if (!(await this.Update())){
                this.$emit('refreshUserInfo')
            }else{
                Vue.prototype.$NotifySuccess("Settings updated")
            }
            this.MailSettingsLoading = false
        },
        Update: async function(): Promise<boolean>{
            //Fix Typings
            this.UserInfo.Settings.Notification.NotificationFrequency = Number.parseInt(this.UserInfo.Settings.Notification.NotificationFrequency)

            const result = await Vue.prototype.$Request("PATCH", "/api/v1/profile/settings", this.UserInfo.Settings)
            if (result.Error == undefined){
                return true
            }
            return false
        }
    }
})
</script>