<template>
  <v-app>
    <v-app-bar
      app
      color="white"
      style="z-index: 100"
    >
      <v-img
        alt="BotAny"
        class="shrink hidden-sm-and-down mr-2"
        src="./assets/logo.png"
        contain
        height="48"
        width="174"
        style="cursor: pointer"
        @click="$router.push('/')"
      />

      <v-img
        alt="BotAny"
        class="shrink hidden-md-and-up mr-2"
        src="./assets/logo-s.png"
        contain
        height="48"
        width="78"
        style="cursor: pointer"
        @click="$router.push('/')"
      />

      <v-scroll-x-transition>
        <top-bar v-if="$route.meta.type==='contest'"></top-bar>
      </v-scroll-x-transition>
      <v-spacer></v-spacer>

      <div class="d-flex" v-if="$vuetify.breakpoint.mdAndUp">
        <div v-if="$store.state.id===-1">
          <v-btn text color="primary"
            :to="{path: '/register/login', query: { redirect: $route.path !== '/error' }}"
          >登录
          </v-btn>
          <v-btn text color="secondary"
            :to="{path: '/register/signup', query: { redirect: $route.path !== '/error' }}"
          >注册
          </v-btn>
          <v-btn text color="primary"
            v-if="$route.meta.type==='login'"
            :to="$store.state.redirect"
          >返回
          </v-btn>
        </div>
        <div v-else>
          <div v-if="!quiting" class="d-flex align-center">
            <router-link text
              class="d-flex align-center"
              style="text-decoration: none"
              :to="`/profile/${$store.state.handle}`">
              <v-avatar tile size="48">
                <v-img
                  :src="$axios.defaults.baseURL + '/user/' + $store.state.handle + '/avatar'"
                  style="border-radius: 4px"
                />
              </v-avatar>
              <span class="ml-2 mr-4">
                <div class="font-weight-bold secondary--text"
                  v-if="$store.state.privilege===$consts.privilege.common"
                >{{$store.state.nickname}}</div>
                <div class="font-weight-bold primary--text"
                  v-if="$store.state.privilege===$consts.privilege.organizer"
                >{{$store.state.nickname}}</div>
                <div class="font-weight-bold secondary--text"
                  v-if="$store.state.privilege===$consts.privilege.superuser"
                >{{$store.state.nickname}}</div>
                <div class="subtitle-2 grey--text">@{{$store.state.handle}}</div>
              </span>
            </router-link>
            <v-btn @click="quiting=true" :loading="quitLoading" :disabled="quitLoading" dark>退出</v-btn>
          </div>
          <div v-else class="d-flex align-center" >
            <div class="secondary--text mr-2">确定退出登录？</div>
            <v-btn @click="logout" dark>是</v-btn>
            <v-btn text color="secondary" @click="quiting=false" class="mr-2" dark>否</v-btn>
          </div>
        </div>
      </div>
    </v-app-bar>

    <v-content>
      <router-view/>
    </v-content>

    <v-footer
      color="#555555"
      dark
    >
      <div>Copyright 2020</div>
    </v-footer>
  </v-app>
</template>

<script>
import TopBar from './components/Topbar.vue'
export default {
  name: 'App',
  components: {
    'top-bar': TopBar
  },
  created () {
    this.$axios.get('/whoami').then(res => {
      this.$store.commit('login', res.data)
    }).catch(err => {
      if (err.response.state === 400) {
        this.$store.commit('logout')
      }
    })
  },
  beforeDestroy () {
    this.$store.commit('logout')
  },
  data: () => ({
    menu: false,
    quiting: false,
    quitLoading: false
  }),
  methods: {
    logout () {
      this.$axios.post('/logout').then(() => {
        this.$store.commit('logout')
        this.quiting = false
        if (this.$route.path === '/') {
          window.location.reload()
        } else {
          this.$router.push('/')
        }
      }).catch(() => {
        this.$store.commit('logout')
        this.quiting = false
        if (this.$route.path === '/') {
          window.location.reload()
        } else {
          this.$router.push('/')
        }
      })
    }
  }
}
</script>
