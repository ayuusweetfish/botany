import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import login from '@/components/login'
import signup from '@/components/signup'
import contestlist from '@/components/contestlist'
import contestmain from '@/components/contestmain'
import contestdetail from '@/components/contestdetail'
import submission from '@/components/submission'
import ranklist from '@/components/ranklist'
import matchlist from '@/components/matchlist'
import match from '@/components/match'
import profile from '@/components/profile'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'login',
      component: login,
      meta: {
        navbarType: 'none'
      }
    },
    {
      path: '/signup',
      name: 'signup',
      component: signup,
      meta: {
        navbarType: 'none'
      }
    },
    {
      path: '/contest_list',
      name: 'contest_list',
      component: contestlist,
      meta: {
        navbarType: 'main'
      }
    },
    {
      path: '/contest_main',
      name: 'contest_main',
      component: contestmain,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/contest_detail',
      name: 'contest_detail',
      component: contestdetail,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/submission',
      name: 'submission',
      component: submission,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/ranklist',
      name: 'ranklist',
      component: ranklist,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/matchlist',
      name: 'matchlist',
      component: matchlist,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/match',
      name: 'match',
      component: match,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/profile',
      name: 'profile',
      component: profile,
      meta: {
        navbarType: 'main'
      }
    }
  ]
})
