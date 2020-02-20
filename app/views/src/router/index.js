import Vue from 'vue'
import VueRouter from 'vue-router'
const ContestList = () => import('../views/ContestList.vue')
const ContestMain = () => import('../views/ContestMain.vue')
const Register = () => import('../views/Register.vue')
const Login = () => import('../views/Login.vue')
const Signup = () => import('../views/Signup.vue')
const Profile = () => import('../views/Profile.vue')

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'ContestList',
    component: ContestList,
    meta: {
      title: '主页',
      type: 'main',
      prePage: [],
      stalling: false
    }
  }, {
    path: '/contest',
    name: 'Contest',
    component: ContestMain,
    meta: {
      title: '赛事',
      type: 'main',
      prePage: [],
      stalling: false
    }
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
          prePage: [],
          stalling: false
        }
      }, {
        path: 'signup',
        component: Signup,
        meta: {
          title: '注册',
          type: 'login',
          prePage: [],
          stalling: true
        }
      }
    ]
  }, {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: {
      title: '选手信息',
      type: 'main',
      prePage: [{ path: '/', query: [] }],
      stalling: false
    }
  }
]

const router = new VueRouter({
  routes
})

export default router
