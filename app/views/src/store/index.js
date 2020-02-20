import Vue from 'vue'
import Vuex from 'vuex'
import Persistance from 'vuex-persistedstate'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    // user
    id: -1,
    handle: '',
    privilege: -1,
    nickname: '',
    routeList: [],
    // contest
    cid: '',
    cname: '',
    // state
    redirect: {
      path: '/'
    },
    stall: false
  },
  mutations: {
    login: (state, user) => {
      state.id = user.id
      state.handle = user.handle
      state.privilege = user.privilege
      state.nickname = user.nickname
    },
    logout: (state) => {
      state.id = -1
      state.handle = ''
      state.privilege = -1
      state.nickname = ''
    },
    setContest: (state, contest) => {
      state.cid = contest.id
      state.cname = contest.title
    },
    setStall: (state, stall) => {
      state.stall = stall
    },
    setRedirect: (state, redirect) => {
      state.redirect = redirect
    }
  },
  actions: {
  },
  modules: {
  },
  plugins: [Persistance()]
})
