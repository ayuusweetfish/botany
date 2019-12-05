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
    nickname: '',
    sitename: '',
    routeList: '',
    beforeLog: null
  },
  mutations: {
    setRouteList (state, val) {
      state.routeList = val
    },
    login (state, val) {
      state.id = val.id
      state.handle = val.handle
      state.email = val.email
      state.privilege = val.privilege
      state.nickname = val.nickname
    },
    logout (state) {
      state.id = ''
      state.handle = ''
      state.email = ''
      state.nickname = ''
      state.privilege = -1
      state.routeList = []
    },
    setPageBefore (state, val) {
      state.beforeLog = val
    },
    enterSubSite (state, val) {
      state.sitename = val
    }
  },
  plugins: [createPersistedState()]
})
