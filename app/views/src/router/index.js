import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import login from '@/components/login'
import register from '@/components/register'
import gamelist from '@/components/gamelist'
import gamemain from '@/components/gamemain'
import gamedetail from '@/components/gamedetail'
import coding from '@/components/coding'
import gameranking from '@/components/gameranking'
import gamevss from '@/components/gamevss'
import vsdetail from '@/components/vsdetail'
import personalinfo from '@/components/personalinfo'

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
      path: '/register',
      name: 'register',
      component: register,
      meta: {
        navbarType: 'none'
      }
    },
    {
      path: '/gamelist',
      name: 'gamelist',
      component: gamelist,
      meta: {
        navbarType: 'main'
      }
    },
    {
      path: '/gamemain',
      name: 'gamemain',
      component: gamemain,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/gamedetail',
      name: 'gamemdetail',
      component: gamedetail,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/coding',
      name: 'coding',
      component: coding,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/gameranking',
      name: 'gameranking',
      component: gameranking,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/gamevss',
      name: 'gamevss',
      component: gamevss,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/vsdetail',
      name: 'vsdetail',
      component: vsdetail,
      meta: {
        navbarType: 'game'
      }
    },
    {
      path: '/personalinfo',
      name: 'personalinfo',
      component: personalinfo,
      meta: {
        navbarType: 'main'
      }
    },
  ]
})
