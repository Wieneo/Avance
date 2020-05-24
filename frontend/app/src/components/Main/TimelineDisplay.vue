<template>
    <div style="margin-top: 5px;">
        <v-skeleton-loader :loading="TicketLoading" transition="fade-transition" type="list-item-three-line">
          <p>Timeline</p>
        </v-skeleton-loader>
    </div>
</template>
<script lang="ts">
  import Vue from 'vue'

   export default Vue.extend({
    name: 'TimelineDisplay',
    props: ["CurrentTicket", "TicketLoading"],
  })
</script>