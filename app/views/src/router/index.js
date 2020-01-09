import Vue from 'vue'
import Router from 'vue-router'
const login = () => import('@/components/login')
const signup = () => import('@/components/signup')
const contestlist = () => import('@/components/contestlist')
const contestmain = () => import('@/components/contestmain')
const submission = () => import('@/components/submission')
const ranklist = () => import('@/components/ranklist')
const matchlist = () => import('@/components/matchlist')
const match = () => import('@/components/match')
const profile = () => import('@/components/profile')
const contestcreate = () => import('@/components/contestcreate')
const contestedit = () => import('@/components/contestedit')
const notfound = () => import('@/components/404')
const submissionlist = () => import('@/components/submissionlist')
const submissioninfo = () => import('@/components/submissioninfo')
const contestscript = () => import('@/components/scriptview')
const contestops = () => import('@/components/contestops')
const judge = () => import('@/components/judge')

Vue.use(Router)

const originalPush = Router.prototype.push
Router.prototype.push = function push (location, onResolve, onReject) {
  if (onResolve || onReject) return originalPush.call(this, location, onResolve, onReject)
  return originalPush.call(this, location).catch(err => err)
}

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
      path: '/submission',
      name: 'submission',
      component: submission,
      meta: {
        title: '提交代码',
        navbarType: 'contest',
        prePage: [{ path: '/contest_main', query: ['cid'] }],
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
        prePage: [{ path: '/contest_main', query: ['cid'] }],
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
        prePage: [{ path: '/contest_main', query: ['cid'] }],
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
          { path: '/contest_main', query: ['cid'] },
          { path: '/match_list', query: ['cid'] }
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
        prePage: [{ path: '/', query: [] }],
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
        prePage: [{ path: '/', query: [] }],
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
        prePage: [
          { path: '/contest_main', query: ['cid'] },
          { path: '/contest_ops', query: ['cid'] }
        ],
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
        prePage: [{ path: '/contest_main', query: ['cid'] }],
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
          { path: '/contest_main', query: ['cid'] },
          { path: '/submission_list', query: ['cid'] }
        ],
        stalling: false
      }
    },
    {
      path: '/contest_script',
      name: 'contest_script',
      component: contestscript,
      meta: {
        title: '赛制脚本',
        navbarType: 'contest',
        prePage: [
          { path: '/contest_main', query: ['cid'] },
          { path: '/contest_ops', query: ['cid'] }
        ],
        stalling: false
      }
    },
    {
      path: '/contest_ops',
      name: 'contestops',
      component: contestops,
      meta: {
        title: '比赛操作',
        navbarType: 'contest',
        prePage: [{ path: '/contest_main', query: ['cid'] }],
        stalling: false
      }
    },
    {
      path: '/judge',
      name: 'judge',
      component: judge,
      meta: {
        title: '裁判程序',
        navbarType: 'contest',
        prePage: [
          { path: '/contest_main', query: ['cid'] },
          { path: '/contest_ops', query: ['cid'] }
        ],
        stalling: false
      }
    }
  ]
})
