import Vue from 'vue'
import Vuex from 'vuex'
import createPersistedState from 'vuex-persistedstate'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    id: '',
    handle: '',
    privilege: -1,
    nickname: '',
    routeList: [],
    contestInfo: {},
    afterLogin: null
  },
  mutations: {
    setRouteList (state, val) {
      state.routeList = val
      console.log(val)
    },
    login (state, val) {
      state.id = val.id
      state.handle = val.handle
      state.privilege = val.privilege
      state.nickname = val.nickname
    },
    logout (state) {
      state.id = ''
      state.handle = ''
      state.nickname = ''
      state.privilege = -1
      state.contestInfo = {}
      state.afterLogin = null
    },
    enterSubSite (state, val) {
      state.contestInfo = val
    },
    clearSubSite (state) {
      state.contestInfo = {}
    },
    setAfterLogin (state, val) {
      state.afterLogin = val
    }
  },
  plugins: [createPersistedState()]
})
