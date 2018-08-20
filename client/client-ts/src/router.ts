import Vue from 'vue';
import Router from 'vue-router';
import Login from '@/views/Login.vue';
import SignUp from '@/views/SignUp.vue';
import ReviewsList from '@/views/ReviewsList.vue';
import NewReview from '@/views/NewReview.vue';
import Profile from '@/views/Profile.vue';
import Review from '@/views/Review.vue';

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
      component: ReviewsList,
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
      component: ReviewsList,
      meta: {
        requiresAuth: true,
      },
      props: {
        type: 'incoming',
      },
    },
    {
      path: '/new-review',
      name: 'new-review',
      component: NewReview,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/profile',
      name: 'profile',
      component: Profile,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/review/:id',
      name: 'review',
      component: Review,
      meta: {
        requiresAuth: true,
      },
    },
  ],
});
