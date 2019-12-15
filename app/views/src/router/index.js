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
        title: '登录',
        navbarType: 'none',
        prePage: []
      }
    },
    {
      path: '/signup',
      name: 'signup',
      component: signup,
      meta: {
        title: '注册',
        navbarType: 'none',
        prePage: []
      }
    },
    {
      path: '/',
      name: 'contest_list',
      component: contestlist,
      meta: {
        title: '主页',
        navbarType: 'main',
        prePage: []
      }
    },
    {
      path: '/contest_main',
      name: 'contest_main',
      component: contestmain,
      meta: {
        title: '比赛主页',
        navbarType: 'contest',
        prePage: []
      }
    },
    {
      path: '/contest_detail',
      name: 'contest_detail',
      component: contestdetail,
      meta: {
        title: '参赛指南',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}]
      }
    },
    {
      path: '/submission',
      name: 'submission',
      component: submission,
      meta: {
        title: '提交代码',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}]
      }
    },
    {
      path: '/ranklist',
      name: 'ranklist',
      component: ranklist,
      meta: {
        title: '选手排行',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}]
      }
    },
    {
      path: '/matchlist',
      name: 'matchlist',
      component: matchlist,
      meta: {
        title: '对局列表',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}]
      }
    },
    {
      path: '/match',
      name: 'match',
      component: match,
      meta: {
        title: '对局信息',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}]
      }
    },
    {
      path: '/profile',
      name: 'profile',
      component: profile,
      meta: {
        title: '选手信息',
        navbarType: 'main',
        prePage: [{path: '/', query: []}]
      }
    },
    {
      path: '/contest_create',
      name: 'contest_create',
      component: contestcreate,
      meta: {
        title: '创建比赛',
        navbarType: 'main',
        prePage: [{path: '/', query: []}]
      }
    },
    {
      path: '/contest_edit',
      name: 'contest_edit',
      component: contestedit,
      meta: {
        title: '修改比赛',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}]
      }
    }
  ]
})
