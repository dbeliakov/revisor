// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'

import axios from 'axios'
import VueAxios from 'vue-axios'

Vue.config.productionTip = false
Vue.use(VueAxios, axios)
Vue.axios.defaults.baseURL = 'https://revisor.dbeliakov.ru/api/'
// Vue.axios.defaults.baseURL = 'http://localhost:8090/api/'
Vue.router = router

Vue.use(require('@websanova/vue-auth'), {
  auth: require('@websanova/vue-auth/drivers/auth/bearer.js'),
  http: require('@websanova/vue-auth/drivers/http/axios.1.x.js'),
  router: require('@websanova/vue-auth/drivers/router/vue-router.2.x.js'),
  refreshData: {enabled: false}
})

Vue.axios.interceptors.response.use(
  (response) => {
    return response
  }, (error) => {
    if (error.response.status === 401) {
      Vue.auth.logout({
        redirect: {name: 'Login'}
      })
    } else if (error.response.status === 500) {
      Vue.router.push({name: '500'})
    }
    return Promise.reject(error)
  })

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router: router,
  components: { App },
  template: '<App/>'
})
