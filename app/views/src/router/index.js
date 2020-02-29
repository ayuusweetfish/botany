import Vue from 'vue'
import VueRouter from 'vue-router'
const Contest = () => import('../views/Contest.vue')
const ContestList = () => import('../views/ContestList.vue')
const ContestMain = () => import('../views/ContestMain.vue')
const Register = () => import('../views/Register.vue')
const Login = () => import('../views/Login.vue')
const Signup = () => import('../views/Signup.vue')
const Profile = () => import('../views/Profile.vue')
const SubmissionList = () => import('../views/SubmissionList.vue')
const SubmissionDetail = () => import('../views/SubmissionDetail.vue')
const ContestEdit = () => import('../views/ContestEdit.vue')
const MatchList = () => import('../views/MatchList.vue')
const MatchDetail = () => import('../views/MatchDetail.vue')
const RankList = () => import('../views/RankList.vue')
const CreateSubmission = () => import('../views/CreateSubmission.vue')
const Judge = () => import('../views/Judge.vue')
const Script = () => import('../views/HandleScript.vue')
const CreateContest = () => import('../views/CreateContest.vue')

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'ContestList',
    component: ContestList,
    meta: {
      title: '主页',
      type: 'main',
      stalling: false
    }
  }, {
    path: '/contest/:cid/main',
    name: 'Contest',
    component: ContestMain,
    meta: {
      title: '赛事',
      type: 'main',
      stalling: false
    }
  }, {
    path: '/create',
    name: 'Create',
    component: CreateContest,
    meta: {
      title: '新建赛事',
      type: 'main',
      stalling: true
    }
  }, {
    path: '/contest/:cid',
    name: 'ContestInfo',
    component: Contest,
    children: [
      {
        path: 'submission/:sid',
        name: 'SubmissionDetail',
        component: SubmissionDetail,
        meta: {
          title: '提交详情',
          type: 'contest',
          stalling: false
        }
      }, {
        path: 'submission',
        name: 'SubmissionList',
        component: SubmissionList,
        meta: {
          title: '提交列表',
          type: 'contest',
          prePage: [{ path: '/contest/main', query: ['cid'] }],
          stalling: false
        }
      }, {
        path: 'edit',
        name: 'ContestEdit',
        component: ContestEdit,
        meta: {
          title: '编辑比赛',
          type: 'contest',
          stalling: true
        }
      }, {
        path: 'match/:mid',
        name: 'MatchDetail',
        component: MatchDetail,
        meta: {
          title: '对局详情',
          type: 'contest',
          stalling: false
        }
      }, {
        path: 'match',
        name: 'MatchList',
        component: MatchList,
        meta: {
          title: '对局列表',
          type: 'contest',
          stalling: false
        }
      }, {
        path: 'ranklist',
        name: 'RankList',
        component: RankList,
        meta: {
          title: '选手排行',
          type: 'contest',
          stalling: false
        }
      }, {
        path: 'participant',
        name: 'Participant',
        component: CreateSubmission,
        meta: {
          title: '我的提交',
          type: 'contest',
          stalling: true
        }
      }, {
        path: 'judge',
        name: 'Judge',
        component: Judge,
        meta: {
          title: '设置裁判',
          type: 'contest',
          stalling: true
        }
      }, {
        path: 'script',
        name: 'Script',
        component: Script,
        meta: {
          title: '脚本操作',
          type: 'contest',
          stalling: false
        }
      }
    ]
  }, {
    path: '/register',
    name: 'Register',
    component: Register,
    children: [
      {
        path: 'login',
        component: Login,
        meta: {
          title: '登录',
          type: 'login',
          stalling: false
        }
      }, {
        path: 'signup',
        component: Signup,
        meta: {
          title: '注册',
          type: 'login',
          stalling: true
        }
      }
    ]
  }, {
    path: '/profile/:handle',
    name: 'Profile',
    component: Profile,
    meta: {
      title: '选手信息',
      type: 'main',
      stalling: false
    }
  }
]

const router = new VueRouter({
  routes
})

export default router
