
<template>
    <v-card>
    <v-navigation-drawer permanent expand-on-hover app
    >
        <v-list>
        <v-list-item class="px-2">
            <v-list-item-avatar>
            <v-img src="https://randomuser.me/api/portraits/women/85.jpg"></v-img>
            </v-list-item-avatar>
        </v-list-item>

        <v-list-item link>
            <v-list-item-content>
            <v-list-item-title class="title">Sandra Adams</v-list-item-title>
            <v-list-item-subtitle>sandra_a88@gmail.com</v-list-item-subtitle>
            </v-list-item-content>
        </v-list-item>
        </v-list>

        <v-divider>
        </v-divider>

        <v-list
        nav
        dense
        >
        <v-list-item link @click="$emit('ShowProjects')">
            <v-list-item-icon>
            <v-icon>mdi-folder</v-icon>
            </v-list-item-icon>
            <v-list-item-title>Projects</v-list-item-title>
        </v-list-item>
        <v-list-item link>
            <v-list-item-icon>
            <v-icon>mdi-star</v-icon>
            </v-list-item-icon>
            <v-list-item-title>Starred</v-list-item-title>
        </v-list-item>
        </v-list>
    </v-navigation-drawer>
    </v-card>
</template>


<script lang="ts">
  import Vue from 'vue'

  export default Vue.extend({
    name: 'Drawer',

    components: {
    }
  })
</script>