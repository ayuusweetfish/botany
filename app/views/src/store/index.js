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
    cname: '',
    myrole: -1,
    redirect: {
      path: '/'
    },
    cstart: new Date(),
    cend: new Date(),
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
      state.cname = contest.title
      state.myrole = contest.my_role
      state.cstart = new Date(contest.start_time * 1000)
      state.cend = new Date(contest.end_time * 1000)
    },
    resetContest: (state) => {
      state.cname = ''
      state.myrole = -1
      state.cstart = new Date()
      state.cend = state.cstart
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
