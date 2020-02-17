import Vue from 'vue'
import VueRouter from 'vue-router'
const ContestList = () => import('../views/ContestList.vue')
const ContestMain = () => import('../views/ContestMain.vue')

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'ContestList',
    component: ContestList
  }, {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  }, {
    path: '/contest',
    name: 'Contest',
    component: ContestMain
  }
]

const router = new VueRouter({
  routes
})

export default router
