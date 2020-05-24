<template>
    <div style="margin-top: 5px;">
        <v-card v-for="action in CurrentTicket.Actions" :key="action.ID">
          <v-card-title>{{action.Title}}</v-card-title>
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