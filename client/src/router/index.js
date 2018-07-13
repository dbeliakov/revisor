import Vue from 'vue'
import Router from 'vue-router'
import Reviews from '@/components/Reviews'
import Login from '@/components/pages/Login'
import SignUp from '@/components/pages/SignUp'
import Profile from '@/components/Profile'
import NewReview from '@/components/pages/NewReview'
import Review from '@/components/Review'
import NotFound from '@/components/pages/NotFound'

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
      path: '/review/:id',
      component: Review,
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
    },
    {
      path: '/404',
      name: 'NotFound',
      component: NotFound
    }
  ]
})
