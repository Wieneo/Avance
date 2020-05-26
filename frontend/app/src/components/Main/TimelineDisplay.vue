<template>
    <div style="margin-top: 5px;">
        <v-skeleton-loader :loading="TicketLoading" transition="fade-transition" type="list-item-three-line" v-if="TicketLoading">
        </v-skeleton-loader>

        <v-timeline v-else dense>
            <v-timeline-item v-for="(action) in CurrentTicket.Actions" :key="action.ID">
              <v-card class="elevation-2">
                <v-card-title class="headline">{{action.IssuedAt | moment("dddd, MM/DD/YYYY HH:mm:ss")}}</v-card-title>
                <v-card-subtitle>{{action.Title}} <i>by</i> {{action.IssuedBy.Firstname}} {{action.IssuedBy.Lastname}} ({{action.IssuedBy.Username}})</v-card-subtitle>
                <v-card-text v-html="action.Content" style="color: black;"></v-card-text>
              </v-card>
            </v-timeline-item>
        </v-timeline>
    </div>
</template>
<script lang="ts">
  import Vue from 'vue'

   export default Vue.extend({
    name: 'TimelineDisplay',
    props: ["CurrentTicket", "TicketLoading"],
  })
</script>