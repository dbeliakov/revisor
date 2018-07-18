import Vue from 'vue';
import Router from 'vue-router';
import Login from '@/views/Login.vue';
import SignUp from '@/views/SignUp.vue';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: {
        requiresNoAuth: true,
      },
    },
    {
      path: '/signup',
      name: 'signup',
      component: SignUp,
      meta: {
        requiresNoAuth: true,
      },
    },
    /*{
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/about',
      name: 'about',
      component: About,
    },*/
  ],
});
