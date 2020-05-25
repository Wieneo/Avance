<template>
    <div style="margin: 5px;">
        <v-skeleton-loader
          class="mx-auto"
          type="text@1, paragraph@3, text@2, paragraph@1"
          v-if="TicketLoading"
        ></v-skeleton-loader>
          <v-card v-for="action in CurrentTicket.Actions" :key="action.ID" style="margin-top: 5px;" v-else>
            <div v-if="action.Type == 0 || action.Type == 1">
              <v-card-subtitle><b>{{action.Title}}</b><br>{{action.IssuedBy.Firstname}} {{action.IssuedBy.Lastname}} ({{action.IssuedBy.Username}})<br>{{action.IssuedAt}}</v-card-subtitle>
              <v-card-text style="color: black;" v-html="action.Content"></v-card-text>
            </div>
          </v-card>
    </div>
</template>
<script lang="ts">
  import Vue from 'vue'

   export default Vue.extend({
    name: 'ActionDisplay',
    props: ["CurrentTicket", "TicketLoading"],
  })
</script>