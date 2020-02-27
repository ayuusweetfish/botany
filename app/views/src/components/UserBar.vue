<template>
  <div>
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
              <div class="font-weight-bold success--text"
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
    <div class="d-flex" v-else>
      <div v-if="$store.state.id===-1">
        <v-btn text color="primary"
          v-if="$route.meta.type!=='login'"
          :to="{path: '/register/login', query: { redirect: $route.path !== '/error' }}"
        >登录
        </v-btn>
        <v-btn text color="secondary"
          v-if="$route.meta.type!=='login'"
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
        <v-menu
          v-model="menu"
          offset-y
          transition="slide-y-transition"
          :close-on-content-click="false"
          min-width="280px"
        >
          <template v-slot:activator="{ on }">
            <v-avatar tile size="48" v-on="on" style="cursor: pointer">
              <v-img
                :src="$axios.defaults.baseURL + '/user/' + $store.state.handle + '/avatar'"
                style="border-radius: 4px"
              />
            </v-avatar>
          </template>
          <v-list>
            <v-list-item @click="goProfile">
              <v-list-item-icon><v-icon>mdi-account-outline</v-icon></v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title
                  class="font-weight-bold secondary--text"
                  v-if="$store.state.privilege===$consts.privilege.common"
                >{{$store.state.nickname}}</v-list-item-title>
                <v-list-item-title
                  class="font-weight-bold primary--text"
                  v-if="$store.state.privilege===$consts.privilege.organizer"
                >{{$store.state.nickname}}</v-list-item-title>
                <v-list-item-title
                  class="font-weight-bold success--text"
                  v-if="$store.state.privilege===$consts.privilege.superuser"
                >{{$store.state.nickname}}</v-list-item-title>
                <v-list-item-subtitle class="subtitle-2 grey--text">@{{$store.state.handle}}</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-btn block color="secondary" v-if="!quiting" @click="quiting=true">退出登录</v-btn>
              <div class="d-flex justify-end align-center" v-else>
                <div>确定退出登录？</div>
                <v-btn color="secondary" @click="logout" :loading="quitLoading" :disabled="quitLoading">
                  确定
                </v-btn>
                <v-btn text color="secondary" @click="quiting=false">
                  取消
                </v-btn>
              </div>
            </v-list-item>
          </v-list>
        </v-menu>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data: () => ({
    menu: false,
    quiting: false,
    quitLoading: false
  }),
  methods: {
    goProfile () {
      this.$router.push(`/profile/${this.$store.state.handle}`)
      this.menu = false
    },
    logout () {
      this.$axios.post('/logout').then(() => {
        this.$store.commit('logout')
        this.quiting = false
        this.menu = false
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

<style>

</style>
