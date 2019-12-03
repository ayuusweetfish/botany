import Vue from 'vue'
import Vuex from 'vuex'
import createPersistedState from 'vuex-persistedstate'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    id: '',
    handle: '',
    email: '',
    privilege: -1,
    nickname: ''
  },
  mutations: {
    setRouteList (state, val) {
      state.routeList = val
    },
    login (state, val) {
      state.userid = val.userid
      state.usertype = val.usertype
      state.username = val.username
    },
    logout (state) {
      state.username = ''
      state.routeList = []
      state.usertype = ''
      state.userid = ''
    }
  },
  plugins: [createPersistedState()]
})