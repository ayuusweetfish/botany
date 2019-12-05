// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
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

Vue.prototype.$qs = qs

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
