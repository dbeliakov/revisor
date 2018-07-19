import Vue from 'vue';
import Router from 'vue-router';
import Login from '@/views/Login.vue';
import SignUp from '@/views/SignUp.vue';
import Reviews from '@/views/Reviews.vue';

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
    { path: '/', redirect: '/outgoing', name: 'home' },
    {
      path: '/outgoing',
      name: 'outgoing',
      component: Reviews,
      meta: {
        requiresAuth: true,
      },
      props: {
        type: 'outgoing',
      },
    },
    {
      path: '/incoming',
      name: 'incoming',
      component: Reviews,
      meta: {
        requiresAuth: true,
      },
      props: {
        type: 'incoming',
      },
    },
  ],
});
