<template>
  <v-container>
    <v-row>
      <v-col>
          <v-card>
          <v-card-title>
            Worker
            <v-spacer></v-spacer>
            <v-text-field
            v-model="CurrentSearch"
            append-icon="mdi-magnify"
            label="Search"
            single-line
            hide-details
          ></v-text-field>
          </v-card-title>
        <v-data-table
          v-model="WorkersSelected"
          :headers="headers"
          :items="Workers"
          item-key="ID"
          show-select
          :loading="Loading"
          :search="CurrentSearch"
          class="elevation-1"
        >
          <template
            v-slot:item.LastSeen="{ item }"
          >{{ item.LastSeen | moment("MM/DD/YYYY HH:mm:ss") }}</template>
          <template v-slot:item.Active="{ item }">
            <v-btn
              icon
              @click="ToggleActive(item)"
              :disabled="!Permissions.Admin && !Permissions.CanChangeWorker"
            >
              <v-icon v-if="item.Active" title="Active">mdi-check-circle-outline</v-icon>
              <v-icon v-else title="Inactive">mdi-bed</v-icon>
            </v-btn>
          </template>
        </v-data-table>
          </v-card>
        <v-btn
          color="green"
          class="mr-2"
          style="margin-top: 10px;"
          :disabled="!Permissions.Admin && !Permissions.CanChangeWorker"
          @click="ToggleSelected(true)"
        >Toggle Active</v-btn>
        <v-btn
          color="orange"
          class="mr-2"
          style="margin-top: 10px;"
          :disabled="!Permissions.Admin && !Permissions.CanChangeWorker"
          @click="ToggleSelected(false)"
        >Toggle Inactive</v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>
<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "Deployment",
  props: ["Permissions"],
  mounted: async function() {
    try {
      if (this.Permissions.Admin || this.Permissions.CanSeeWorker) {
        this.Workers = await Vue.prototype.$Request("GET", "/api/v1/workers");
      } else {
        this.WorkerPermsFailed = true;
      }
    } finally {
      this.Loading = false;
    }
  },
  data: function() {
    return {
      CurrentSearch: '',
      Loading: true,
      Workers: [],
      WorkerPermsFailed: false,
      WorkersSelected: [],
      headers: [
        { text: "ID", align: "start", value: "ID" },
        { text: "Name", value: "Name" },
        { text: "Last Seen at", value: "LastSeen" },
        { text: "Active", value: "Active" }
      ]
    };
  },
  methods: {
    ToggleActive: async function(Data: any) {
      //Returns the new worker state
      this.Loading = true;
      try {
        const result = await Vue.prototype.$Request(
          "PATCH",
          "/api/v1/worker/" + Data.ID
        );
        this.Workers.forEach(element => {
          if (element.ID == Data.ID) {
            element.Active = result.Active;
          }
        });
      } finally {
        this.Loading = false;
      }
    },
    ToggleSelected: async function(Enable: boolean){
      this.Loading = true
      try{
        this.asyncForEach(this.WorkersSelected, async (element) => {
          if (Enable != element.Active){
            const result = await Vue.prototype.$Request(
              "PATCH",
              "/api/v1/worker/" + element.ID
            );

            element.Active = result.Active
          }
        });
      }finally{
        this.Loading = false
      }
    },
    asyncForEach: async function (array: any, callback: any) {
        for (let index = 0; index < array.length; index++) {
          await callback(array[index], index, array);
        }
      }
  }
});
</script>