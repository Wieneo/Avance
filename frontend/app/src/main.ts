import Vue from 'vue'
import App from './App.vue'
import store from './store'
import vuetify from './plugins/vuetify';
import VueRouter from 'vue-router'
import MainApp from './components/Main/MainApp.vue'
import Login from './components/Login/Login.vue'
import Settings from './components/Settings/Settings.vue'
import {Utils} from './plugins/utils'

import Moment from 'vue-moment'

Vue.config.productionTip = false

Vue.use(VueRouter)
Vue.use(Utils)
Vue.use(Moment);

const router = new VueRouter({
  mode: 'history',
  routes: [
    // dynamic segments start with a colon
    { path: '/', component: MainApp, name: "Main" },
    { path: '/login', component: Login, name: "Login" },
    { path: '/settings', component: Settings, name: "Settings" }
  ]
})

new Vue({
  store,
  vuetify,
  router,
  render: h => h(App)
}).$mount('#app')
