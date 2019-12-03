// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'


Vue.config.productionTip = false
Vue.use(ElementUI)
<<<<<<< HEAD
=======

Vue.prototype.$axios = axios
axios.defaults.baseURL = '/api'


>>>>>>> parent of 20eb236... Merge branch 'frontend' of github.com:kawa-yoiko/botany into backend-dev
/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
