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
import notfound from '@/components/404'
import submissionlist from '@/components/submissionlist'
import submissioninfo from '@/components/submissioninfo'
import contestscript from '@/components/scriptview'

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
        prePage: [],
        stalling: false
      }
    },
    {
      path: '/signup',
      name: 'signup',
      component: signup,
      meta: {
        title: '注册',
        navbarType: 'none',
        prePage: [],
        stalling: false
      }
    },
    {
      path: '/',
      name: 'contest_list',
      component: contestlist,
      meta: {
        title: '主页',
        navbarType: 'main',
        prePage: [],
        stalling: false
      }
    },
    {
      path: '/contest_main',
      name: 'contest_main',
      component: contestmain,
      meta: {
        title: '比赛主页',
        navbarType: 'contest',
        prePage: [],
        stalling: false
      }
    },
    {
      path: '/contest_detail',
      name: 'contest_detail',
      component: contestdetail,
      meta: {
        title: '参赛指南',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}],
        stalling: false
      }
    },
    {
      path: '/submission',
      name: 'submission',
      component: submission,
      meta: {
        title: '提交代码',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}],
        stalling: false
      }
    },
    {
      path: '/ranklist',
      name: 'ranklist',
      component: ranklist,
      meta: {
        title: '选手排行',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}],
        stalling: false
      }
    },
    {
      path: '/match_list',
      name: 'match_list',
      component: matchlist,
      meta: {
        title: '对局列表',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}],
        stalling: false
      }
    },
    {
      path: '/match',
      name: 'match',
      component: match,
      meta: {
        title: '对局信息',
        navbarType: 'contest',
        prePage: [
          {path: '/contest_main', query: ['cid']},
          {path: '/match_list', query: ['cid']}
        ],
        stalling: false
      }
    },
    {
      path: '/profile',
      name: 'profile',
      component: profile,
      meta: {
        title: '选手信息',
        navbarType: 'main',
        prePage: [{path: '/', query: []}],
        stalling: false
      }
    },
    {
      path: '/contest_create',
      name: 'contest_create',
      component: contestcreate,
      meta: {
        title: '创建比赛',
        navbarType: 'main',
        prePage: [{path: '/', query: []}],
        stalling: true
      }
    },
    {
      path: '/contest_edit',
      name: 'contest_edit',
      component: contestedit,
      meta: {
        title: '修改比赛',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}],
        stalling: true
      }
    },
    {
      path: '/notfound',
      name: 'notfound',
      component: notfound,
      meta: {
        title: '出错了',
        navbarType: 'main',
        prePage: [],
        stalling: false
      }
    },
    {
      path: '/submission_list',
      name: 'submission_list',
      component: submissionlist,
      meta: {
        title: '提交列表',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}],
        stalling: false
      }
    },
    {
      path: '/submission_info',
      name: 'submission_info',
      component: submissioninfo,
      meta: {
        title: '提交详情',
        navbarType: 'contest',
        prePage: [
          {path: '/contest_main', query: ['cid']},
          {path: '/submission_list', query: ['cid']}
        ],
        stalling: false
      }
    },
    {
      path: '/contest_script',
      name: 'contest_script',
      component: contestscript,
      meta: {
        title: '比赛脚本',
        navbarType: 'contest',
        prePage: [{path: '/contest_main', query: ['cid']}],
        stalling: false
      }
    }
  ]
})
