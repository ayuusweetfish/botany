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
import contestcreate from '@/components/contestcreate'
import contestedit from '@/components/contestedit'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/login',
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
      path: '/',
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
        navbarType: 'contest'
      }
    },
    {
      path: '/contest_detail',
      name: 'contest_detail',
      component: contestdetail,
      meta: {
        navbarType: 'contest'
      }
    },
    {
      path: '/submission',
      name: 'submission',
      component: submission,
      meta: {
        navbarType: 'contest'
      }
    },
    {
      path: '/ranklist',
      name: 'ranklist',
      component: ranklist,
      meta: {
        navbarType: 'contest'
      }
    },
    {
      path: '/matchlist',
      name: 'matchlist',
      component: matchlist,
      meta: {
        navbarType: 'contest'
      }
    },
    {
      path: '/match',
      name: 'match',
      component: match,
      meta: {
        navbarType: 'contest'
      }
    },
    {
      path: '/profile',
      name: 'profile',
      component: profile,
      meta: {
        navbarType: 'main'
      }
    },
    {
      path: '/contest_create',
      name: 'contest_create',
      component: contestcreate,
      meta: {
        navbarType: 'main'
      }
    },
    {
      path: '/contest_edit',
      name: 'contest_edit',
      component: contestedit,
      meta: {
        navbarType: 'contest'
      }
    }
  ]
})
