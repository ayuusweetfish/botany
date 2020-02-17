import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import Axios from 'axios'
import utils from './utils'
import qs from 'querystring'
import LoadingOverlay from './components/LoadingOverlay.vue'

Vue.config.productionTip = false
Axios.defaults.baseURL = '/api'
Axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'
Vue.prototype.$axios = Axios
Vue.prototype.$qs = qs
Vue.use(utils)
Vue.component('v-loading-overlay', LoadingOverlay)

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
