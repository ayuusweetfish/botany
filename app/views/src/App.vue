<template>
  <div id="app">
    <el-row class = "topbar" type="flex" justify="space-between">
      <el-col :span="8" v-if="$route.meta.navbarType !== 'contest'" class="topbar-tittle">
        <div align='left' @click="$router.push('/')" style="cursor: pointer">Botany-Demo</div>
      </el-col>
      <el-col :span="8" v-if="$route.meta.navbarType === 'contest'" class="topbar-tittle">
        <div style="display: inline">{{$store.state.contestInfo.title}}</div>
        <div style="display: inline">@Botany</div>
      </el-col>
      <!-- <el-col :span="6" v-if="$route.meta.navbarType === 'main'" align='left'>
        <div style="display: inline; color:gray">|</div>
        <router-link class="navbar-item" to="/">比赛列表</router-link>
        <div style="display: inline; color:gray">|</div>
        <router-link class="navbar-item" to="/profile">个人信息</router-link>
        <div style="display: inline; color:gray">|</div>
      </el-col> -->
      <el-col :span="10" v-if="$route.meta.navbarType === 'contest'" align='center'>
        <router-link class="navbar-item" :to="{path:'/contest_main', query: {cid: $route.query.cid}}">比赛首页</router-link>
        <router-link class="navbar-item" :to="{path:'/contest_detail', query: {cid: $route.query.cid}}">参赛指南</router-link>
        <router-link v-if="checkRouteValid('imIn')" class="navbar-item" :to="{path: '/submission', query: {cid: $route.query.cid}}">我的代码</router-link>
        <!-- <el-button type="text" class="navbar-item" @click="gocontestranking">查看排行</el-button>
        <div style="display: inline; color:gray">|</div>
        <el-button type="text" class="navbar-item" @click="gocontestvss">查看对局</el-button>
        <div style="display: inline; color:gray">|</div> -->
        <router-link class="navbar-item" to="/">返回</router-link>
      </el-col>
      <el-col :span="4" v-if="$route.meta.navbarType !== 'none'">
        <el-dropdown v-if="$store.state.handle" :hide-on-click="true" @command="handleCommand" trigger="click">
          <span class="el-dropdown-link" style="cursor: pointer;">
            <el-avatar :src="defaultAva"></el-avatar>
          </span>
          <el-dropdown-menu  slot="dropdown" style="min-width: 120px; padding: 2px 10px 2px 10px;">
            <el-dropdown-item :disabled="true" class="info-dropdown-item" style="color: #505050; font-weight: 600; border-bottom: 1px solid silver">{{$store.state.nickname}}</el-dropdown-item>
            <!-- <el-dropdown-item :disabled="true" class="info-dropdown-item">UID: {{$store.state.id}}</el-dropdown-item> -->
            <el-dropdown-item :disabled="true" class="info-dropdown-item" style="color: grey">{{translatePrivilege($store.state.privilege)}}</el-dropdown-item>
            <el-dropdown-item :disabled="true" class="info-dropdown-item">账号：{{$store.state.handle}}</el-dropdown-item>
            <el-dropdown-item :disabled="true" class="info-dropdown-item">UID：{{$store.state.id}}</el-dropdown-item>
            <router-link :to="{path: '/profile', query: {handle: $store.state.handle}}" style="text-decoration: none">
              <el-dropdown-item class="button-dropdown-item" style="border-top: 1px solid silver">我的资料</el-dropdown-item>
            </router-link>
            <el-dropdown-item command="logout" class="button-dropdown-item">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
        <div v-else>
          <el-link :underline="false" type="primary" class="login-button" @click="goLogin">登录</el-link>
          <el-divider direction="vertical"></el-divider>
          <el-link :underline="false" type="" class="login-button">注册</el-link>
        </div>
      </el-col>
    </el-row>
    <el-row style="margin-bottom: 10px">
      <el-col :span="24">
        <div v-if="$route.meta.navbarType !== 'none'" align="left">
          <i class = "el-icon-caret-right" v-if="$store.state.routeList.length !== 0" style="display: inline-block"></i>
          <el-breadcrumb separator="/" style = "display: inline-block; margin-left: 2px">
            <el-breadcrumb-item v-for="(item) in $store.state.routeList" :key="item.path">
              <router-link :to="{path: item.path, query: item.query}">{{item.title}}</router-link>
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
      </el-col>
    </el-row>
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App',
  created () {
    this.$axios.get('/whoami').then(res => {
      console.log(res.data)
      this.$store.commit('login', {
        id: res.data.id,
        handle: res.data.handle,
        privilege: res.data.privilege,
        nickname: res.data.nickname
      })
    }).catch(err => {
      if (err.response.state === 400) {
        this.$store.commit('logout')
      }
    })
  },
  beforeDestroy () {
    this.$store.commit('logout')
  },
  data () {
    return {
      showPwdDlg: false,
      showPhnDlg: false,
      defaultAva: require('./assets/logo.png')
    }
  },
  methods: {
    changeTitle (title) {
      console.log(title)
    },
    goLogin () {
      console.log(this.$route)
      this.$store.commit('setAfterLogin', {path: this.$route.path, query: this.$route.query})
      this.$router.push({
        path: '/login',
        query: {
          redirect: true
        }
      })
    },
    handleCommand (cmd) {
      if (cmd === 'logout') {
        this.logout()
      } else if (cmd === 'toProfile') {
        this.toProfile()
      }
    },
    translatePrivilege (num) {
      let str = ''
      switch (num) {
        case this.$consts.privilege.common:
          str = 'Common User'
          break
        case this.$consts.privilege.oranizer:
          str = 'Organizer'
          break
        case this.$consts.privilege.superuser:
          str = 'Super User'
          break
        default:
          break
      }
      return str
    },
    logout () {
      this.$axios.post('/logout').then(res => {
        this.$store.commit('logout')
        window.location.reload()
      }).catch(err => {
        console.log(err)
      })
    },
    checkRouteValid (role) {
      if (this.$store.state.contestInfo &&
      this.$store.state.contestInfo.my_role &&
      this.$store.state.contestInfo.my_role === this.$consts.role[role]) {
        return true
      } else {
        return false
      }
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin: auto;
  margin-top: 0px;
  width: 1080px;
  font-family:'Segoe UI', Tahoma, Geneva, Verdana, sans-serif
}
.topbar {
  border-bottom: 1px solid silver;
  max-width: 1080px;
  margin: auto;
  align-items: flex-end;
  margin-bottom: 10px;
  min-width: 720px;
  min-height: 84px;
}
.topbar-tittle {
  font-weight: 600;
  font-size: 30px;
  text-align: left;
  margin-left: 20px;
  line-height: 60px;
}
.navbar-item {
  font-size: 18px;
  font-weight: 500;
  color:#545454;
  line-height: 48px;
  text-decoration: none;
  margin: 0px 10px 0px 10px;
}
.info-dropdown-item {
  font-size: 14px;
  padding: 0
}
.button-dropdown-item {
  font-size: 14px;
  padding: 0;
}
.login-button {
  font-size: 18px;
  font-weight: 600;
  height: 30px;
}
</style>
