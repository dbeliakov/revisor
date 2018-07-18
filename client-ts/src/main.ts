import Vue from 'vue';
import App from './App.vue';
import router from '@/router';
import Axios from 'axios';
import {AuthPlugin} from '@/auth/vue-plugin';

const axios = Axios.create({
  baseURL: 'http://localhost:8090/api',
});

Vue.config.productionTip = false;

Vue.use(AuthPlugin, {axios, router});

new Vue({
  router,
  render: (h) => h(App),
}).$mount('#app');
