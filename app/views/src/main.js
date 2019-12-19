// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store'
import ElementUI from 'element-ui'
import axios from 'axios'
import qs from 'querystring'
import utils from './utils'
import 'element-ui/lib/theme-chalk/index.css'

Vue.config.productionTip = false
Vue.use(ElementUI)
Vue.use(utils)

Vue.prototype.$axios = axios
axios.defaults.baseURL = '/api'
axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'
axios.interceptors.response.use(
  function (res) {
    return res
  },
  function (error) {
    if (error.response.status === 401) {
      ElementUI.MessageBox.confirm('你需要先登录', '提示', {
        type: 'warning',
        confirmButtonText: '登录'
      }).then(res => {
        store.commit('setAfterLogin', {
          path: router.currentRoute.path,
          query: router.currentRoute.query
        })
        router.push({
          path: '/login',
          query: {
            redirect: true
          }
        }).catch()
      })
    } else if (error.response.status === 403) {
      ElementUI.Message.error('没有操作权限')
    } else if (error.response.status === 404) {
      ElementUI.Message.error('404 Not Found')
    }
    return Promise.reject(error)
  }
)

Vue.prototype.$qs = qs

router.beforeEach((to, from, next) => {
  let routeList = []
  if (!to.meta.prePage) {
    return next('/notfound')
  }
  if (from.meta.stalling && store.state.stallFlag) {
    const check = window.confirm('表单尚未提交，确定离开?')
    if (!check) {
      return next(false)
    }
  }
  to.meta.prePage.forEach(item => {
    let r = router.resolve(item.path).route
    let page = {
      title: r.meta.title,
      path: item.path,
      query: {}
    }
    item.query.forEach(key => {
      page.query[key] = to.query[key]
    })
    routeList.push(page)
  })
  routeList.push({
    title: to.meta.title,
    path: to.path,
    query: to.query
  })
  store.commit('setRouteList', routeList)
  next()
})

window.onbeforeunload = function () {
  if (router.currentRoute.meta.stalling && store.state.stallFlag) {
    return ''
  }
}

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
