import Vue from 'vue'
import Router from 'vue-router'
import Reviews from '@/components/Reviews'
import Login from '@/components/Login'
import SignUp from '@/components/SignUp'
import Profile from '@/components/Profile'
import NewReview from '@/components/NewReview'

Vue.use(Router)

export default new Router({
  routes: [
    { path: '/', redirect: '/outgoing', name: 'Home' },
    {
      path: '/outgoing',
      name: 'OutReviews',
      component: Reviews,
      meta: {
        auth: true
      },
      props: {
        type: 'outgoing'
      }
    },
    {
      path: '/incoming',
      name: 'InReviews',
      component: Reviews,
      meta: {
        auth: true
      },
      props: {
        type: 'incoming'
      }
    },
    {
      path: '/review/new',
      name: 'NewReview',
      component: NewReview,
      meta: {
        auth: true
      }
    },
    {
      path: '/login',
      name: 'Login',
      component: Login,
      meta: {
        auth: false
      }
    },
    {
      path: '/signup',
      name: 'SignUp',
      component: SignUp,
      meta: {
        auth: false
      }
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile,
      meta: {
        auth: true
      }
    }
  ]
})
