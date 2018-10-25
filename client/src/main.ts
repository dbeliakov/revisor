import Vue from 'vue';
import App from './App.vue';
import router from '@/router';
import Axios from 'axios';
import {AuthPlugin} from '@/auth/vue-plugin';
import {ReviewsPlugin} from '@/reviews/vue-plugin';

const axios = Axios.create({
  baseURL: process.env.VUE_APP_API_URL,
});

console.log(process.env)

Vue.config.productionTip = false;

Vue.use(AuthPlugin, {axios, router});
Vue.use(ReviewsPlugin, {axios});

new Vue({
  router,
  render: (h) => h(App),
}).$mount('#app');
